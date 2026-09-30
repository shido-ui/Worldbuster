package api

import("encoding/json";"net/http";"github.com/shido-ui/Worldbuster/backend/internal/social")
type socialHandler struct{s *social.Service}
func newSocialHandler(s *social.Service)*socialHandler{return &socialHandler{s:s}}
func(h *socialHandler)send(w http.ResponseWriter,r *http.Request){var q struct{FromID string `json:"fromId"`;ToID string `json:"toId"`;Body string `json:"body"`};if json.NewDecoder(r.Body).Decode(&q)!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return};m,err:=h.s.Send(q.FromID,q.ToID,q.Body);if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,201,m)}
func(h *socialHandler)inbox(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");writeJSON(w,200,h.s.Inbox(id))}
