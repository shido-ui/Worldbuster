package api

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"

 "github.com/shido-ui/Worldbuster/backend/internal/auth"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
 "github.com/shido-ui/Worldbuster/backend/internal/events"
 "github.com/shido-ui/Worldbuster/backend/internal/inventory"
 "github.com/shido-ui/Worldbuster/backend/internal/job"
 "github.com/shido-ui/Worldbuster/backend/internal/organization"
 "github.com/shido-ui/Worldbuster/backend/internal/progression"
 "github.com/shido-ui/Worldbuster/backend/internal/simulation"
 "github.com/shido-ui/Worldbuster/backend/internal/social"
 "github.com/shido-ui/Worldbuster/backend/internal/world"
)

func TestSimulationProductionPathThroughREST(t *testing.T) {
 ws := world.NewService()
 jobs := job.NewService()
 if err := jobs.Register(job.Job{ID:"general-work",Name:"General Workforce",Positions:[]job.Position{{ID:"worker",JobID:"general-work",Name:"Worker"}}}); err != nil {
  t.Fatal(err)
 }
 if err := jobs.Employ("npc-rest-e2e","general-work","worker",0,0); err != nil {
  t.Fatal(err)
 }

 population := simulation.NewService()
 if err := population.Register(simulation.SimCharacter{
  ID:"npc-rest-e2e",Name:"Worker",Controller:simulation.ControllerSimulated,Active:true,
  Goals:[]simulation.Goal{{Kind:"WORK",Priority:10}},
 }); err != nil {
  t.Fatal(err)
 }
 state := simulation.NewStateService()
 eventsService := events.NewService()
 adapter := &simulation.LiveAdapter{
  Jobs:jobs,
  Progression:progression.NewService(),
  Social:social.NewService(),
  Travel:world.NewTravelService(),
  Locations:map[string]string{"npc-rest-e2e":"central"},
 }
 runner := &simulation.IntegratedRunner{
  Population:population,
  State:state,
  World:adapter,
  Emit:func(eventType,actorID,targetID string,payload map[string]any){
   eventsService.Publish(eventType,actorID,targetID,payload)
  },
 }
 runtime := simulation.NewRuntime(runner,simulation.TierScheduler{Policy:simulation.LifecyclePolicy{ActivePercent:100,RecentPercent:0,BackgroundPercent:0}})
 if got := runtime.Tick(time.Unix(1,0)); got != 1 {
  t.Fatalf("production simulation runner executed %d actions, want 1",got)
 }

 authBackend := MemoryAuthBackend{Service:auth.NewService()}
 router := NewRouter(ws,authBackend,world.NewTravelService(),inventory.NewService(),economy.NewService(),organization.NewService(),jobs,progression.NewService(),social.NewService(),eventsService)
 req := httptest.NewRequest(http.MethodGet,"/api/v1/events?limit=10",nil)
 rr := httptest.NewRecorder()
 router.Handler().ServeHTTP(rr,req)
 if rr.Code != http.StatusOK {
  t.Fatalf("events REST endpoint returned %d, want %d",rr.Code,http.StatusOK)
 }

 var got []map[string]any
 if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
  t.Fatalf("decode REST response: %v",err)
 }
 found := false
 for _,event := range got {
  if event["type"] == "SIMULATED_ACTION" && event["actorId"] == "npc-rest-e2e" {
   found = true
   break
  }
 }
 if !found {
  t.Fatalf("REST response did not contain the production simulation event: %+v",got)
 }
}
