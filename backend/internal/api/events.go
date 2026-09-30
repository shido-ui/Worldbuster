package api

import("net/http";"strconv";"github.com/shido-ui/Worldbuster/backend/internal/events")
type eventHandler struct{s *events.Service}
func newEventHandler(s *events.Service)*eventHandler{return &eventHandler{s:s}}
func(h *eventHandler)recent(w http.ResponseWriter,r *http.Request){limit:=50;if n,e:=strconv.Atoi(r.URL.Query().Get("limit"));e==nil&&n>0{limit=n};writeJSON(w,200,h.s.Recent(limit))}
func(h *eventHandler)notifications(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");if id==""{writeJSON(w,400,map[string]string{"error":"characterId is required"});return};writeJSON(w,200,h.s.Notifications(id))}
