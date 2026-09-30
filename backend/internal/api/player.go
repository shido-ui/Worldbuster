package api

import (
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/player"
)

func (r *Router) playerProfile(w http.ResponseWriter,req *http.Request){
 if r.playerContext==nil {writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=r.playerContext.Resolve(req.Context(),req)
 if err!=nil {writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 writeJSON(w,http.StatusOK,p)
}

var _ player.Profile
