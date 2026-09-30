package api

import (
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
 "github.com/shido-ui/Worldbuster/backend/internal/inventory"
 "github.com/shido-ui/Worldbuster/backend/internal/player"
)

func (r *Router) playerProfile(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil {writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=r.playerContext.Resolve(req.Context(),req)
 if err!=nil {writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 writeJSON(w,http.StatusOK,p)
}

var _ player.Profile

func(r *Router) playerDashboard(w http.ResponseWriter,req *http.Request){if r.playerContext==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return};session,p,err:=r.playerContext.Resolve(req.Context(),req);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return};e,err:=r.playerContext.Economy.GetByAccountID(req.Context(),session.AccountID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"economy unavailable"});return};inv,err:=r.playerContext.Inventory.Get(req.Context(),p.ID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"inventory unavailable"});return};writeJSON(w,http.StatusOK,PlayerDashboard{Profile:p,Economy:e,Inventory:inv})}
