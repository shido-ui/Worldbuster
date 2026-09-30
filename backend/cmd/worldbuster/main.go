package main

import(
 "log"
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
)

func main(){
 ws:=world.NewService();as:=auth.NewService();ts:=world.NewTravelService();is:=inventory.NewService();es:=economy.NewService();os:=organization.NewService();js:=job.NewService();ps:=progression.NewService();ss:=social.NewService();evs:=events.NewService()
 _=is.RegisterItem(inventory.Item{ID:"water",Name:"Water",Category:"supply",Stackable:true,MaxStack:10})

 population:=simulation.NewService()
 state:=simulation.NewStateService()
 adapter:=&simulation.LiveAdapter{Jobs:js,Progression:ps,Social:ss,Travel:ts,Locations:map[string]string{}}
 runner:=&simulation.IntegratedRunner{Population:population,State:state,World:adapter,Emit:func(t,actor,target string,payload map[string]any){evs.Publish(t,actor,target,payload)}}

 go func(){ticker:=time.NewTicker(time.Second);defer ticker.Stop();for now:=range ticker.C{ws.Tick();runner.Tick(now)}}()
 server:=&http.Server{Addr:":8080",Handler:api.NewRouter(ws,as,ts,is,es,os,js,ps,ss,evs).Handler(),ReadHeaderTimeout:5*time.Second}
 log.Println("Worldbuster server listening on :8080");log.Fatal(server.ListenAndServe())
}
