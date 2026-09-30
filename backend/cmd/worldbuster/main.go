package main

import(
 "log"
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
 ws:=world.NewService();var as api.AuthBackend;if dbStore.SQL!=nil { as=store.NewDatabaseAuthService(dbStore) } else { as=api.MemoryAuthBackend{Service:auth.NewService()} }ts:=world.NewTravelService();is:=inventory.NewService();es:=economy.NewService();os:=organization.NewService();js:=job.NewService();ps:=progression.NewService();ss:=social.NewService();evs:=events.NewService()
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
 for _,c:=range generator.Generate(50){
  if err:=population.Register(c);err!=nil{log.Printf("simulation population: %v",err);continue}
  locations[c.ID]="central"
  _=js.Employ(c.ID,"general-work","worker",0,0)
 }
 runner:=&simulation.IntegratedRunner{Population:population,State:state,World:adapter,Emit:func(t,actor,target string,payload map[string]any){evs.Publish(t,actor,target,payload)}}
 runtime:=simulation.NewRuntime(runner,simulation.TierScheduler{Policy:simulation.LifecyclePolicy{ActivePercent:20,RecentPercent:30,BackgroundPercent:30}})

 go func(){ticker:=time.NewTicker(time.Second);defer ticker.Stop();for now:=range ticker.C{ws.Tick();runtime.Tick(now)}}()
 router:=api.NewRouter(ws,as,ts,is,es,os,js,ps,ss,evs)
 if dbStore.SQL!=nil {
  pr:=store.PlayerRepository{DB:dbStore}
  router.WithProgressionRepository(&api.ProgressionAPI{Repo:store.ProgressionRepository{DB:dbStore},ResolvePlayer:func(ctx context.Context,accountID string)(string,error){p,err:=pr.GetByAccountID(ctx,accountID);return p.ID,err}})
  router.WithPlayerContext(api.NewPlayerContext(dbStore,as))
  router.WithPersistentJobs(&api.PersistentJobsAPI{Repo:store.JobRepository{DB:dbStore},Context:api.NewPlayerContext(dbStore,as)})
 }
 server:=&http.Server{Addr:":8080",Handler:router.Handler(),ReadHeaderTimeout:5*time.Second}
 log.Println("Worldbuster server listening on :8080");log.Fatal(server.ListenAndServe())
}
