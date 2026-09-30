package api

import (
 "encoding/json"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type MissionsAPI struct{Repo store.MissionRepository;Context *PlayerContext}

func(a *MissionsAPI) List(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 x,err:=a.Repo.List(r.Context(),p.Level);if err!=nil{writeJSON(w,500,map[string]string{"error":"missions unavailable"});return};active,err:=a.Repo.ListForPlayer(r.Context(),p.ID);if err!=nil{writeJSON(w,500,map[string]string{"error":"player missions unavailable"});return};writeJSON(w,200,map[string]any{"missions":x,"playerMissions":active})
}
func(a *MissionsAPI) Accept(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 var q struct{MissionID string `json:"missionId"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.MissionID==""{writeJSON(w,400,map[string]string{"error":"missionId required"});return}
 x,err:=a.Repo.Accept(r.Context(),q.MissionID,p.ID,p.Level);if err!=nil{writeJSON(w,400,map[string]string{"error":"mission cannot be accepted"});return};writeJSON(w,201,x)
}
func(a *MissionsAPI) Progress(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 var q struct{MissionID string `json:"missionId"`;Amount int64 `json:"amount"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.MissionID==""||q.Amount<=0{writeJSON(w,400,map[string]string{"error":"missionId and positive amount required"});return}
 x,err:=a.Repo.Progress(r.Context(),q.MissionID,p.ID,q.Amount);if err!=nil{writeJSON(w,400,map[string]string{"error":"mission progress rejected"});return};writeJSON(w,200,x)
}
func(a *MissionsAPI) Complete(w http.ResponseWriter,r *http.Request){
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,401,map[string]string{"error":"authentication required"});return}
 var q struct{MissionID string `json:"missionId"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.MissionID==""{writeJSON(w,400,map[string]string{"error":"missionId required"});return}
 x,e,err:=a.Repo.Complete(r.Context(),q.MissionID,p.ID);if err!=nil{writeJSON(w,400,map[string]string{"error":"mission cannot be completed"});return};writeJSON(w,200,map[string]any{"mission":x,"reward":e})
}
