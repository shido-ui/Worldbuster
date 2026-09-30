package main

import("log";"net/http";"time";"github.com/shido-ui/Worldbuster/backend/internal/api";"github.com/shido-ui/Worldbuster/backend/internal/auth";"github.com/shido-ui/Worldbuster/backend/internal/inventory";"github.com/shido-ui/Worldbuster/backend/internal/economy";"github.com/shido-ui/Worldbuster/backend/internal/organization";"github.com/shido-ui/Worldbuster/backend/internal/job";"github.com/shido-ui/Worldbuster/backend/internal/progression";"github.com/shido-ui/Worldbuster/backend/internal/social";"github.com/shido-ui/Worldbuster/backend/internal/world")
func main(){
 ws:=world.NewService();as:=auth.NewService();ts:=world.NewTravelService();is:=inventory.NewService();es:=economy.NewService();os:=organization.NewService();js:=job.NewService();ps:=progression.NewService();ss:=social.NewService();_=is.RegisterItem(inventory.Item{ID:"water",Name:"Water",Category:"supply",Stackable:true,MaxStack:10})
 go func(){ticker:=time.NewTicker(time.Second);defer ticker.Stop();for range ticker.C{ws.Tick()}}()
 server:=&http.Server{Addr:":8080",Handler:api.NewRouter(ws,as,ts,is,es,os,js,ps,ss).Handler(),ReadHeaderTimeout:5*time.Second}
 log.Println("Worldbuster server listening on :8080");log.Fatal(server.ListenAndServe())
}
