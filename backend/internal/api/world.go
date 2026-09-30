package api

import("encoding/json";"net/http";"time";"github.com/shido-ui/Worldbuster/backend/internal/world")

type worldHandler struct{travelService *world.TravelService; context *PlayerContext}
func newWorldHandler(t *world.TravelService,pc *PlayerContext)*worldHandler{return &worldHandler{travelService:t,context:pc}}
func(h *worldHandler)locations(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]any{"locations":world.Locations(),"routes":world.Routes()})}
func(h *worldHandler)currentTravel(w http.ResponseWriter,r *http.Request){
 if h.context==nil{writeJSON(w,http.StatusServiceUnavailable,map[string]string{"error":"player persistence unavailable"});return}
 _,p,err:=h.context.Resolve(r.Context(),r);if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
 now:=time.Now().UTC();t,ok:=h.travelService.Get(p.ID,now);if ok&&!now.Before(t.ArrivesAt){if err:=h.context.Profile.SetLocation(r.Context(),p.ID,t.To);err!=nil{writeJSON(w,http.StatusInternalServerError,map[string]string{"error":"location persistence failed"});return};p.LocationID=t.To;ok=false};writeJSON(w,http.StatusOK,map[string]any{"traveling":ok,"travel":func()any{if !ok{return nil};return t}(),"locationId":p.LocationID})
}
func(h *worldHandler)travel(w http.ResponseWriter,r *http.Request){
 if r.Method!="POST"{w.WriteHeader(http.StatusMethodNotAllowed);return}
 var req struct{CharacterID string `json:"characterId"`;From string `json:"from"`;To string `json:"to"`}
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return}
 if req.CharacterID==""||req.From==""||req.To==""{writeJSON(w,400,map[string]string{"error":"characterId, from and to are required"});return}
 if h.context!=nil{
  _,p,err:=h.context.Resolve(r.Context(),r)
  if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
  if p.ID!=req.CharacterID{writeJSON(w,http.StatusForbidden,map[string]string{"error":"character does not belong to authenticated player"});return};from:=p.LocationID;if from==""{from="central"};if req.From!=from{writeJSON(w,http.StatusConflict,map[string]string{"error":"travel origin does not match authoritative player location"});return};req.From=from
 }
 t,err:=h.travelService.Start(req.CharacterID,req.From,req.To,time.Now().UTC())
 if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return}
 writeJSON(w,202,t)
}
