package api

import (
 "encoding/json"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type ProgressionAPI struct{ Repo store.ProgressionRepository }

func(a ProgressionAPI) AddXP(w http.ResponseWriter,r *http.Request){
 if r.Method!="POST"{http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
 var req struct{PlayerID string `json:"playerId"`;Amount int64 `json:"amount"`}
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil||req.PlayerID==""||req.Amount<=0{http.Error(w,"invalid request",http.StatusBadRequest);return}
 result,err:=a.Repo.AddXPAndUnlocks(r.Context(),req.PlayerID,req.Amount,nil)
 if err!=nil{http.Error(w,"progression update failed",http.StatusInternalServerError);return}
 w.Header().Set("Content-Type","application/json")
 _=json.NewEncoder(w).Encode(result)
}
