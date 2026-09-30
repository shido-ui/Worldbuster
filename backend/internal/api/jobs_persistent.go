package api

import (
 "encoding/json"
 "errors"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type PersistentJobsAPI struct{Repo store.JobRepository; Context *PlayerContext}

func(a *PersistentJobsAPI) List(w http.ResponseWriter,r *http.Request){jobs,err:=a.Repo.List(r.Context());if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"jobs unavailable"});return};writeJSON(w,http.StatusOK,map[string]any{"jobs":jobs})}

func(a *PersistentJobsAPI) Employ(w http.ResponseWriter,r *http.Request){
 if a.Context==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=a.Context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 var req map[string]string;if json.NewDecoder(r.Body).Decode(&req)!=nil||req["jobId"]==""{writeJSON(w,http.StatusBadRequest,map[string]string{"error":"jobId required"});return}
 stat:=p.Strength+p.Defense+p.Speed+p.Intelligence+p.Endurance
 j,err:=a.Repo.Employ(r.Context(),p.ID,req["jobId"],p.Level,stat)
 if errors.Is(err,store.ErrAlreadyEmployed){writeJSON(w,http.StatusConflict,map[string]string{"error":"already employed"});return}
 if errors.Is(err,store.ErrJobNotFound){writeJSON(w,http.StatusNotFound,map[string]string{"error":"job not found"});return}
 if errors.Is(err,store.ErrJobRequirement){writeJSON(w,http.StatusForbidden,map[string]string{"error":"job requirements not met"});return}
 if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"employment failed"});return}
 writeJSON(w,http.StatusOK,map[string]any{"employed":true,"job":j})
}