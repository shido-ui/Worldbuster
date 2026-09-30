package api

import (
 "net/http"
 "strconv"
 "github.com/shido-ui/Worldbuster/backend/internal/events"
)

type eventHandler struct{s *events.Service; context *PlayerContext}
func newEventHandler(s *events.Service,pc *PlayerContext)*eventHandler{return &eventHandler{s:s,context:pc}}

func(h *eventHandler)recent(w http.ResponseWriter,r *http.Request){
 limit:=50
 if n,e:=strconv.Atoi(r.URL.Query().Get("limit"));e==nil&&n>0{limit=n}
 if limit>100{limit=100}
 writeJSON(w,http.StatusOK,h.s.Recent(limit))
}

func(h *eventHandler)notifications(w http.ResponseWriter,r *http.Request){
 if h.context==nil{writeJSON(w,http.StatusForbidden,map[string]string{"error":"notifications require player context"});return}
 _,profile,err:=h.context.Resolve(r.Context(),r)
 if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 if profile.ID==""{writeJSON(w,http.StatusNotFound,map[string]string{"error":"player profile not found"});return}
 writeJSON(w,http.StatusOK,h.s.Notifications(profile.ID))
}
