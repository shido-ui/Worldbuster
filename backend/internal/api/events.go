package api

import (
 "net/http"
 "encoding/json"
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

func(h *eventHandler)markNotificationRead(w http.ResponseWriter,r *http.Request){
 if h.context==nil{writeJSON(w,http.StatusForbidden,map[string]string{"error":"notifications require player context"});return}
 _,profile,err:=h.context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 id:=r.URL.Query().Get("id");if id==""{writeJSON(w,http.StatusBadRequest,map[string]string{"error":"notification id required"});return}
 if !h.s.MarkRead(profile.ID,id){writeJSON(w,http.StatusNotFound,map[string]string{"error":"notification not found"});return};writeJSON(w,http.StatusOK,map[string]any{"read":true})
}
func(h *eventHandler)markNotificationsRead(w http.ResponseWriter,r *http.Request){
 if h.context==nil{writeJSON(w,http.StatusForbidden,map[string]string{"error":"notifications require player context"});return}
 _,profile,err:=h.context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return};writeJSON(w,http.StatusOK,map[string]any{"marked":h.s.MarkAllRead(profile.ID)})
}

func(h *eventHandler)stream(w http.ResponseWriter,r *http.Request){
 if h.context==nil{http.Error(w,"authentication required",401);return};if _,_,err:=h.context.Resolve(r.Context(),r);err!=nil{http.Error(w,"authentication required",401);return}
 f,ok:=w.(http.Flusher);if !ok{http.Error(w,"streaming unsupported",500);return};w.Header().Set("Content-Type","text/event-stream");w.Header().Set("Cache-Control","no-cache");w.Header().Set("Connection","keep-alive");ch:=h.s.Subscribe();defer h.s.Unsubscribe(ch);f.Flush();
 for{select{case <-r.Context().Done():return;case ev,ok:=<-ch:if !ok{return};b,err:=json.Marshal(ev);if err!=nil{continue};if _,err=w.Write([]byte("event: world\ndata: "));err!=nil{return};if _,err=w.Write(b);err!=nil{return};if _,err=w.Write([]byte("\n\n"));err!=nil{return};f.Flush()}}}
