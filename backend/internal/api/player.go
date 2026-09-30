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

func(r *Router) playerEconomy(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 session,p,err:=r.playerContext.Resolve(req.Context(),req);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 account,err:=r.playerContext.Economy.GetByAccountID(req.Context(),p.AccountID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"economy unavailable"});return}
 ledger,err:=r.playerContext.Economy.LedgerByAccountID(req.Context(),session.AccountID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"ledger unavailable"});return}
 writeJSON(w,http.StatusOK,map[string]any{"account":account,"ledger":ledger})
}

func(r *Router) playerState(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 session,p,err:=r.playerContext.Resolve(req.Context(),req);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 e,err:=r.playerContext.Economy.GetByAccountID(req.Context(),session.AccountID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"economy unavailable"});return}
 inv,err:=r.playerContext.Inventory.Get(req.Context(),p.ID);if err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"inventory unavailable"});return}
 writeJSON(w,http.StatusOK,map[string]any{"profile":p,"economy":e,"inventory":inv})
}
