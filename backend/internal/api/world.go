package api

import("encoding/json";"net/http";"time";"github.com/shido-ui/Worldbuster/backend/internal/world")

type worldHandler struct{travelService *world.TravelService; context *PlayerContext}
func newWorldHandler(t *world.TravelService,pc *PlayerContext)*worldHandler{return &worldHandler{travelService:t,context:pc}}
func(h *worldHandler)locations(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]any{"locations":world.Locations(),"routes":world.Routes()})}
func(h *worldHandler)travel(w http.ResponseWriter,r *http.Request){
 if r.Method!="POST"{w.WriteHeader(http.StatusMethodNotAllowed);return}
 var req struct{CharacterID string `json:"characterId"`;From string `json:"from"`;To string `json:"to"`}
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return}
 if req.CharacterID==""||req.From==""||req.To==""{writeJSON(w,400,map[string]string{"error":"characterId, from and to are required"});return}
 if h.context!=nil{
  _,p,err:=h.context.Resolve(r.Context(),r)
  if err!=nil{writeJSON(w,http.StatusUnauthorized,map[string]string{"error":"authentication required"});return}
  if p.ID!=req.CharacterID{writeJSON(w,http.StatusForbidden,map[string]string{"error":"character does not belong to authenticated player"});return}
 }
 t,err:=h.travelService.Start(req.CharacterID,req.From,req.To,time.Now().UTC())
 if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return}
 writeJSON(w,202,t)
}
