package main

import(
 "log"
 "fmt"
 "database/sql"
 "context"
 "os"
 _ "github.com/lib/pq"
 "net/http"
 "time"
 "github.com/shido-ui/Worldbuster/backend/internal/api"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
 "github.com/shido-ui/Worldbuster/backend/internal/inventory"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
 "github.com/shido-ui/Worldbuster/backend/internal/organization"
 "github.com/shido-ui/Worldbuster/backend/internal/job"
 "github.com/shido-ui/Worldbuster/backend/internal/progression"
 "github.com/shido-ui/Worldbuster/backend/internal/social"
 "github.com/shido-ui/Worldbuster/backend/internal/events"
 "github.com/shido-ui/Worldbuster/backend/internal/world"
 "github.com/shido-ui/Worldbuster/backend/internal/simulation"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

func main(){
 dsn:=os.Getenv("WORLDBUSTER_DATABASE_URL")
 var db *sql.DB
 if dsn!="" { var err error; db,err=sql.Open("postgres",dsn); if err!=nil {log.Fatal(err)}; defer db.Close(); if err=db.Ping(); err!=nil {log.Fatal(err)} }
 dbStore:=store.New(db)
 if dbStore.SQL!=nil { if err:=store.ApplyMigrations(context.Background(),dbStore,"database/migrations"); err!=nil { log.Fatal(err) } }
 ws:=world.NewService();var as api.AuthBackend;if dbStore.SQL!=nil { as=store.NewDatabaseAuthService(dbStore) } else { as=api.MemoryAuthBackend{Service:auth.NewService()} }
 ts:=world.NewTravelService();is:=inventory.NewService();es:=economy.NewService();os:=organization.NewService();js:=job.NewService();ps:=progression.NewService();ss:=social.NewService();evs:=events.NewService()
 _=is.RegisterItem(inventory.Item{ID:"water",Name:"Water",Category:"supply",Stackable:true,MaxStack:10})

 // Minimal deterministic seed content keeps the simulation runnable while the full data catalog is built.
 _=js.Register(job.Job{ID:"general-work",Name:"General Workforce",Department:"Operations",Positions:[]job.Position{{ID:"worker",JobID:"general-work",Name:"Worker",Level:1,BaseSalary:100,RequiredEducation:0,RequiredStat:0}}})
 _=ps.RegisterCourse(progression.Course{ID:"orientation",Name:"World Orientation",DurationHours:1,EducationGain:1,RequiredLevel:1})

 population:=simulation.NewService()
 state:=simulation.NewStateService()
 locations:=map[string]string{}
 adapter:=&simulation.LiveAdapter{Jobs:js,Progression:ps,Social:ss,Travel:ts,Locations:locations}
 generator:=simulation.NewPopulationGenerator(42,simulation.PopulationProfile{
  Names:[]string{"Aster","Vale","Rin","Kade","Mira","Nox","Sora","Iris"},
 })
 generated:=generator.Generate(50)
 for _,c:=range generated{
  if err:=population.Register(c);err!=nil{log.Printf("simulation population: %v",err);continue}
  locations[c.ID]="central"
  _=js.Employ(c.ID,"general-work","worker",0,0)
  if dbStore.SQL!=nil {
   simRepo:=store.SimulationRepository{DB:dbStore}
   if err:=simRepo.EnsureCharacter(context.Background(),c);err!=nil {log.Printf("simulation persistence: %v",err)}
  }
 }
 runner:=&simulation.IntegratedRunner{Population:population,State:state,World:adapter,Emit:func(t,actor,target string,payload map[string]any){evs.Publish(t,actor,target,payload)}}
 if dbStore.SQL!=nil {
  prodRepo:=store.ProductionRepository{DB:dbStore}
  for i,c:=range generated {
   if i>=5 {break}
   bType:=simulation.BusinessProduction
   b:=simulation.BusinessState{ID:simulation.NewBusinessID(),OwnerID:c.ID,Name:c.Name+" Works",Type:bType,Cash:500,Level:1,Active:true}
   if b.ID=="" {continue}
   if err:=prodRepoCreate(dbStore, b);err!=nil {log.Printf("npc business seed: %v",err);continue}
   input:="raw-water";qty:=25
   if i%2==1 {input="fiber"}
   if err:=prodRepo.SeedBusinessInventory(context.Background(),b.ID,input,int64(qty));err!=nil {log.Printf("npc supply seed: %v",err)}
  }
 }
 if dbStore.SQL!=nil { runner.Persistence=store.SimulationRepository{DB:dbStore} }
 runtime:=simulation.NewRuntime(runner,simulation.TierScheduler{Policy:simulation.LifecyclePolicy{ActivePercent:20,RecentPercent:30,BackgroundPercent:30}})

 go func(){ticker:=time.NewTicker(time.Second);defer ticker.Stop();for now:=range ticker.C{ws.Tick();runtime.Tick(now)}}()
 if dbStore.SQL!=nil {
  go func(){ticker:=time.NewTicker(time.Minute);defer ticker.Stop();for range ticker.C{
   result,err:=store.ProductionRepository{DB:dbStore}.RunCycle(context.Background(),100)
   if err!=nil {log.Printf("production cycle: %v",err)} else if result.Produced>0 {evs.Publish("world.production.cycle","world","world",map[string]any{"businessesProcessed":result.Processed,"produced":result.Produced,"inputsConsumed":result.InputsConsumed,"outputsCreated":result.OutputsCreated,"laborSpent":result.LaborSpent})}
   territory,tberr:=store.TerritoryRepository{DB:dbStore}.ResolveConsequences(context.Background(),250)
   if tberr!=nil {log.Printf("territory consequences: %v",tberr)} else if territory.ControlChanges>0 || territory.StabilityChanges>0 {evs.Publish("world.territory.cycle","world","world",map[string]any{"territories":territory.TerritoriesProcessed,"controlChanges":territory.ControlChanges,"stabilityChanges":territory.StabilityChanges})}
   balance,berr:=store.EconomyBalanceRepository{DB:dbStore}.Rebalance(context.Background(),250)
   if berr==nil && balance.PriceChanges>0 { if we,err:=store.WorldEventRepository{DB:dbStore}.Create(context.Background(),"market-fluctuation",1,"central",map[string]any{"priceChanges":balance.PriceChanges}); err==nil { evs.Publish("world.event.created","world",we.ID,map[string]any{"code":we.Code,"severity":we.Severity}); if _,nerr:=store.NewsRepository{DB:dbStore}.Publish(context.Background(),we.ID,"Market conditions shift","Local market prices changed across "+fmt.Sprint(balance.PriceChanges)+" tracked assets.","ECONOMY",1,"central"); nerr!=nil {log.Printf("news publish: %v",nerr)} } }
   if berr!=nil {log.Printf("economy rebalance: %v",berr)} else if balance.PriceChanges>0 {evs.Publish("world.economy.rebalanced","world","world",map[string]any{"assets":balance.ProcessedAssets,"priceChanges":balance.PriceChanges,"supply":balance.TotalSupply,"demand":balance.TotalDemand})}

   payments,err:=store.JobRepository{DB:dbStore}.SettleDueSalaries(context.Background(),60)
   if err!=nil {log.Printf("salary settlement: %v",err);continue}
   for _,pay:=range payments {evs.Publish("job.salary.paid",pay.PlayerID,pay.PlayerID,map[string]any{"amount":pay.Amount,"jobId":pay.JobID,"xp":25,"level":pay.Level})}
  }}()
 }
 router:=api.NewRouter(ws,as,ts,is,es,os,js,ps,ss,evs)
if dbStore.SQL!=nil {
  router.WithEducation(&api.EducationAPI{Repo:store.EducationRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})

  pr:=store.PlayerRepository{DB:dbStore}
  router.WithProgressionRepository(&api.ProgressionAPI{Repo:store.ProgressionRepository{DB:dbStore},ResolvePlayer:func(ctx context.Context,accountID string)(string,error){p,err:=pr.GetByAccountID(ctx,accountID);return p.ID,err}})
  router.WithPlayerContext(api.NewPlayerContext(dbStore,as))
  router.WithPersistentJobs(&api.PersistentJobsAPI{Repo:store.JobRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithPersistentSocial(&api.PersistentSocialAPI{Repo:store.SocialRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithReputation(&api.ReputationAPI{Repo:store.SocialRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithPersistentOrganizations(&api.PersistentOrganizationsAPI{Repo:store.OrganizationRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithTerritory(&api.TerritoryAPI{Repo:store.TerritoryRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithMarket(&api.MarketAPI{Repo:store.MarketRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithMissions(&api.MissionsAPI{Repo:store.MissionRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithCombat(&api.CombatAPI{Repo:store.CombatRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithEquipment(&api.EquipmentAPI{Repo:store.EquipmentRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
  router.WithWorldEvents(&api.WorldEventsAPI{Repo:store.WorldEventRepository{DB:dbStore}})
  router.WithNews(&api.NewsAPI{Repo:store.NewsRepository{DB:dbStore}})
 }
  if dbStore.SQL!=nil {
  go func(){ticker:=time.NewTicker(time.Minute);defer ticker.Stop();for range ticker.C{
   if _,err:=(store.EducationRepository{DB:dbStore}).CompleteDue(context.Background());err!=nil{log.Printf("education completion: %v",err)}
  }}()
 }
 server:=&http.Server{Addr:":8080",Handler:router.Handler(),ReadHeaderTimeout:5*time.Second}
 log.Println("Worldbuster server listening on :8080");log.Fatal(server.ListenAndServe())
}

func prodRepoCreate(dbStore *store.DB,b simulation.BusinessState) error {
 return (store.SimulationRepository{DB:dbStore}).CreateBusiness(context.Background(),b)
}
