package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/shido-ui/Worldbuster/backend/internal/api"
	"github.com/shido-ui/Worldbuster/backend/internal/cache"
	"github.com/shido-ui/Worldbuster/backend/internal/economy"
	"github.com/shido-ui/Worldbuster/backend/internal/events"
	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/organization"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/scheduler"
	"github.com/shido-ui/Worldbuster/backend/internal/seed"
	"github.com/shido-ui/Worldbuster/backend/internal/simulation"
	"github.com/shido-ui/Worldbuster/backend/internal/social"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

const defaultAddress = ":8080"

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, closeDB, err := openDatabase(ctx)
	if err != nil {
		logger.Error("database startup failed", "error", err)
		return
	}
	defer closeDB()

	deps, runtime, err := buildDependencies(ctx, db, logger)
	if err != nil {
		logger.Error("application startup failed", "error", err)
		return
	}

	router := ConfigureRouter(deps)
	jobs := buildSchedulerJobs(deps, runtime, logger)
	schedulerRunner := scheduler.New(logger)
	schedulerRunner.Start(ctx, jobs...)

	server := &http.Server{
		Addr:              configuredAddress(),
		Handler:           router.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("Worldbuster server listening", "addr", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server stopped", "error", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP graceful shutdown failed", "error", err)
		}
	}

	stop()
	schedulerRunner.Wait()
}

func configuredAddress() string {
	if address := os.Getenv("WORLDBUSTER_ADDR"); address != "" {
		return address
	}
	return defaultAddress
}

func openDatabase(ctx context.Context) (*store.DB, func(), error) {
	dsn := os.Getenv("WORLDBUSTER_DATABASE_URL")
	if dsn == "" {
		return nil, func() {}, errors.New("WORLDBUSTER_DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, func() {}, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, func() {}, err
	}

	dbStore := store.New(db)
	migrationDir := os.Getenv("WORLDBUSTER_MIGRATIONS_DIR")
	if migrationDir == "" {
		migrationDir = "database/migrations"
	}
	if err := store.ApplyMigrations(ctx, dbStore, migrationDir); err != nil {
		_ = db.Close()
		return nil, func() {}, err
	}

	return dbStore, func() { _ = db.Close() }, nil
}

type simulationRuntime struct {
	world      *world.Service
	runtime    *simulation.Runtime
	events     *events.Service
	population *simulation.Service
}

func buildDependencies(ctx context.Context, db *store.DB, logger *slog.Logger) (Dependencies, *simulationRuntime, error) {
	ws := world.NewService()
	travel := world.NewTravelService()
	inventoryService := inventory.NewService()
	economyService := economy.NewService()
	organizationService := organization.NewService()
	jobService := job.NewService()
	progressionService := progression.NewService()
	socialService := social.NewService()
	eventService := events.NewService()

	redisClient, err := buildRedisClient(ctx)
	if err != nil {
		return Dependencies{}, nil, err
	}
	seedDir := os.Getenv("WORLDBUSTER_SEED_DIR")
	content, err := seed.LoadContent(seedDir)
	if err != nil {
		return Dependencies{}, nil, err
	}
	for _, item := range content.Items {
		if err := inventoryService.RegisterItem(item); err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed item %s: %w", item.ID, err)
		}
	}
	for _, seededJob := range content.Jobs {
		if err := jobService.Register(seededJob); err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed job %s: %w", seededJob.ID, err)
		}
	}
	for _, course := range content.Courses {
		if err := progressionService.RegisterCourse(course); err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed course %s: %w", course.ID, err)
		}
	}

	authBackend := api.AuthBackend(store.NewDatabaseAuthService(db))

	population := simulation.NewService()
	state := simulation.NewStateService()
	locations := map[string]string{}
	adapter := &simulation.LiveAdapter{
		Jobs: jobService, Progression: progressionService, Social: socialService,
		Travel: travel, Locations: locations,
	}

	generator := simulation.NewPopulationGenerator(content.Population.Seed, simulation.PopulationProfile{
		Names:           content.Population.Names,
		JobWeights:      content.Population.JobWeights,
		PersonalityBias: content.Population.PersonalityBias,
	})
	populationSize := 50
	if raw := os.Getenv("WORLDBUSTER_SIM_POPULATION"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 || n > 10000 {
			logger.Warn("invalid WORLDBUSTER_SIM_POPULATION; using default", "value", raw, "default", populationSize)
		} else {
			populationSize = n
		}
	}

	generated := generator.Generate(populationSize)
	simRepo := store.SimulationRepository{DB: db}
	for _, character := range generated {
		if err := population.Register(character); err != nil {
			return Dependencies{}, nil, fmt.Errorf("register simulation character %s: %w", character.ID, err)
		}
		locations[character.ID] = "central"
		if err := jobService.Employ(character.ID, "general-work", "worker", 0, 0); err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed employment for %s: %w", character.ID, err)
		}
		if err := simRepo.EnsureCharacter(ctx, character); err != nil {
			return Dependencies{}, nil, fmt.Errorf("persist simulation character %s: %w", character.ID, err)
		}
	}

	runner := &simulation.IntegratedRunner{
		Population: population,
		State:      state,
		World:      adapter,
		Emit: func(eventType, actorID, targetID string, payload map[string]any) {
			publishEvent(eventService, eventType, actorID, targetID, payload, logger)
		},
		PersistenceError: func(operation, characterID string, err error) {
			logger.Error("simulation persistence failed", "operation", operation, "character", characterID, "error", err)
		},
	}
	runner.Persistence = simRepo

	runtime := simulation.NewRuntime(
		runner,
		simulation.TierScheduler{
			Policy: simulation.LifecyclePolicy{ActivePercent: 20, RecentPercent: 30, BackgroundPercent: 30},
		},
	)

	prodRepo := store.ProductionRepository{DB: db}
	for i, character := range generated {
		if i >= 5 {
			break
		}
		businessID := simulation.DeterministicBusinessID(character.ID, simulation.BusinessProduction)
		business := simulation.BusinessState{
			ID: businessID, OwnerID: character.ID, Name: character.Name + " Works",
			Type: simulation.BusinessProduction, Cash: 500, Level: 1, Active: true,
		}
		persistedID, err := prodRepoCreate(db, business)
		if err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed NPC business: %w", err)
		}
		input := "raw-water"
		if i%2 == 1 {
			input = "fiber"
		}
		if err := prodRepo.SeedBusinessInventory(ctx, persistedID, input, 25); err != nil {
			return Dependencies{}, nil, fmt.Errorf("seed NPC business inventory: %w", err)
		}
	}

	return Dependencies{
			DB: db, Cache: redisClient, Auth: authBackend, World: ws, Travel: travel,
			Inventory: inventoryService, Economy: economyService,
			Organization: organizationService, Job: jobService,
			Progression: progressionService, Social: socialService, Events: eventService,
		}, &simulationRuntime{
			world: ws, runtime: runtime, events: eventService, population: population,
		}, nil
}

func buildSchedulerJobs(deps Dependencies, runtime *simulationRuntime, logger *slog.Logger) []scheduler.Job {
	jobs := []scheduler.Job{
		{
			Name:     "world-simulation",
			Interval: time.Second,
			Run: func(_ context.Context, now time.Time) error {
				runtime.world.Tick()
				runtime.runtime.Tick(now)
				return nil
			},
		},
	}

	jobs = append(jobs,
		scheduler.Job{
			Name:     "production-cycle",
			Interval: time.Minute,
			Run: func(ctx context.Context, _ time.Time) error {
				result, err := (store.ProductionRepository{DB: deps.DB}).RunCycle(ctx, 100)
				if err != nil {
					return err
				}
				if result.Produced > 0 {
					publishEvent(runtime.events, "world.production.cycle", "world", "world", map[string]any{
						"businessesProcessed": result.Processed, "produced": result.Produced,
						"inputsConsumed": result.InputsConsumed, "outputsCreated": result.OutputsCreated,
						"laborSpent": result.LaborSpent,
					}, logger)
				}
				return nil
			},
		},
		scheduler.Job{
			Name:     "territory-consequences",
			Interval: time.Minute,
			Run: func(ctx context.Context, _ time.Time) error {
				result, err := (store.TerritoryRepository{DB: deps.DB}).ResolveConsequences(ctx, 250)
				if err != nil {
					return err
				}
				if result.ControlChanges > 0 || result.StabilityChanges > 0 {
					publishEvent(runtime.events, "world.territory.cycle", "world", "world", map[string]any{
						"territories":      result.TerritoriesProcessed,
						"controlChanges":   result.ControlChanges,
						"stabilityChanges": result.StabilityChanges,
					}, logger)
				}
				return nil
			},
		},
		scheduler.Job{
			Name:     "economy-rebalance",
			Interval: time.Minute,
			Run: func(ctx context.Context, _ time.Time) error {
				balance, err := (store.EconomyBalanceRepository{DB: deps.DB}).Rebalance(ctx, 250)
				if err != nil {
					return err
				}
				if balance.PriceChanges > 0 {
					worldEventsRepo := store.WorldEventRepository{DB: deps.DB}
					event, createErr := worldEventsRepo.Create(ctx, "market-fluctuation", 1, "central", map[string]any{
						"priceChanges": balance.PriceChanges,
					})
					if createErr != nil {
						return createErr
					}
					publishEvent(runtime.events, "world.event.created", "world", event.ID, map[string]any{
						"code": event.Code, "severity": event.Severity,
					}, logger)
					_, newsErr := (store.NewsRepository{DB: deps.DB}).Publish(
						ctx, event.ID, "Market conditions shift",
						"Local market prices changed across "+fmt.Sprint(balance.PriceChanges)+" tracked assets.",
						"ECONOMY", 1, "central",
					)
					if newsErr != nil {
						return newsErr
					}
					publishEvent(runtime.events, "world.economy.rebalanced", "world", "world", map[string]any{
						"assets": balance.ProcessedAssets, "priceChanges": balance.PriceChanges,
						"supply": balance.TotalSupply, "demand": balance.TotalDemand,
					}, logger)
				}
				return nil
			},
		},
		scheduler.Job{
			Name:     "salary-settlement",
			Interval: time.Minute,
			Run: func(ctx context.Context, _ time.Time) error {
				payments, err := (store.JobRepository{DB: deps.DB}).SettleDueSalaries(ctx, 60)
				if err != nil {
					return err
				}
				for _, payment := range payments {
					publishEvent(runtime.events, "job.salary.paid", payment.PlayerID, payment.PlayerID, map[string]any{
						"amount": payment.Amount, "jobId": payment.JobID, "xp": 25, "level": payment.Level,
					}, logger)
				}
				return nil
			},
		},
		scheduler.Job{
			Name:     "education-completion",
			Interval: time.Minute,
			Run: func(ctx context.Context, _ time.Time) error {
				_, err := (store.EducationRepository{DB: deps.DB}).CompleteDue(ctx)
				return err
			},
		},
	)

	return jobs
}

func prodRepoCreate(db *store.DB, business simulation.BusinessState) (string, error) {
	return (store.SimulationRepository{DB: db}).CreateBusiness(context.Background(), business)
}

func publishEvent(service *events.Service, eventType, actorID, targetID string, payload map[string]any, logger *slog.Logger) {
	if _, err := service.Publish(eventType, actorID, targetID, payload); err != nil {
		logger.Error("event publish failed", "eventType", eventType, "error", err)
	}
}

func buildRedisClient(ctx context.Context) (*cache.Client, error) {
	addr := os.Getenv("WORLDBUSTER_REDIS_ADDR")
	if addr == "" {
		return nil, nil
	}
	password := os.Getenv("WORLDBUSTER_REDIS_PASSWORD")
	client, err := cache.New(addr, password, 0)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis startup failed: %w", err)
	}
	return client, nil
}
