package app

import (
	"context"
	"fmt"

	"github.com/shido-ui/Worldbuster/backend/internal/api"
	"github.com/shido-ui/Worldbuster/backend/internal/cache"
	"github.com/shido-ui/Worldbuster/backend/internal/economy"
	"github.com/shido-ui/Worldbuster/backend/internal/events"
	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/organization"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/social"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

type Dependencies struct {
	Cache        *cache.Client
	DB           *store.DB
	Auth         api.AuthBackend
	World        *world.Service
	Travel       *world.TravelService
	Inventory    *inventory.Service
	Economy      *economy.Service
	Organization *organization.Service
	Job          *job.Service
	Progression  *progression.Service
	Social       *social.Service
	Events       *events.Service
}

func ConfigureRouter(d Dependencies) *api.Router {
	router := api.NewRouter(
		d.World,
		d.Auth,
		d.Travel,
		d.Inventory,
		d.Economy,
		d.Organization,
		d.Job,
		d.Progression,
		d.Social,
		d.Events,
	)

	playerContext := api.NewPlayerContext(d.DB, d.Auth)
	var playerRepo store.PlayerStore = store.PlayerRepository{DB: d.DB}

	router.
		WithEducation(&api.EducationAPI{
			Repo:    store.EducationRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithProgressionRepository(&api.ProgressionAPI{
			Repo: store.ProgressionRepository{DB: d.DB},
			ResolvePlayer: func(ctx context.Context, accountID string) (string, error) {
				player, err := playerRepo.GetByAccountID(ctx, accountID)
				if err != nil {
					return "", err
				}
				return player.ID, nil
			},
		}).
		WithPlayerContext(playerContext).
		WithPersistentJobs(&api.PersistentJobsAPI{
			Repo:    store.JobRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithPersistentSocial(&api.PersistentSocialAPI{
			Repo:    store.SocialRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithReputation(&api.ReputationAPI{
			Repo:    store.SocialRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithPersistentOrganizations(&api.PersistentOrganizationsAPI{
			Repo:    store.OrganizationRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithTerritory(&api.TerritoryAPI{
			Repo:    store.TerritoryRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithMarket(&api.MarketAPI{
			Repo:    store.MarketRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithMissions(&api.MissionsAPI{
			Repo:    store.MissionRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithCombat(&api.CombatAPI{
			Repo:    store.CombatRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithEquipment(&api.EquipmentAPI{
			Repo:    store.EquipmentRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithWorldEvents(&api.WorldEventsAPI{
			Repo: store.WorldEventRepository{DB: d.DB},
		}).
		WithNews(&api.NewsAPI{
			Repo: store.NewsRepository{DB: d.DB},
		}).
		WithAchievements(&api.AchievementsAPI{
			Repo:    store.AchievementRepository{DB: d.DB},
			Context: playerContext,
		}).
		WithAdmin(&api.AdminAPI{
			Repo: store.AdminRepository{DB: d.DB},
			Auth: d.Auth,
		}).
		WithMetrics(api.NewMetrics()).
		WithReadiness(func(ctx context.Context) error {
			if d.DB == nil || d.DB.SQL == nil {
				return fmt.Errorf("database is not configured")
			}
			if err := d.DB.SQL.PingContext(ctx); err != nil {
				return fmt.Errorf("database: %w", err)
			}
			if d.Cache != nil {
				if err := d.Cache.Ping(ctx); err != nil {
					return fmt.Errorf("redis: %w", err)
				}
			}
			return nil
		})
	return router
}
