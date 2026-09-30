package api

import("encoding/json";"net/http";"github.com/shido-ui/Worldbuster/backend/internal/world")

type worldHandler struct{travel *world.TravelService}
func newWorldHandler(t *world.TravelService)*worldHandler{return &worldHandler{travel:t}}
func(h *worldHandler)locations(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]any{"locations":world.Locations(),"routes":world.Routes()})}
func(h *worldHandler)travel(w http.ResponseWriter,r *http.Request){
 if r.Method!="POST"{w.WriteHeader(http.StatusMethodNotAllowed);return}
 var req struct{CharacterID string `json:"characterId"`;From string `json:"from"`;To string `json:"to"`}
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return}
 t,err:=h.travel.Start(req.CharacterID,req.From,req.To,timeNowUTC())
 if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return}
 writeJSON(w,202,t)
}
func timeNowUTC()time.Time{return time.Now().UTC()}
