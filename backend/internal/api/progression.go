package api

import("encoding/json";"net/http";"github.com/shido-ui/Worldbuster/backend/internal/progression")
type progressionHandler struct{s *progression.Service}
func newProgressionHandler(s *progression.Service)*progressionHandler{return &progressionHandler{s:s}}
func(h *progressionHandler)get(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");if id==""{writeJSON(w,400,map[string]string{"error":"characterId is required"});return};writeJSON(w,200,h.s.Get(id))}
func(h *progressionHandler)courses(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,h.s.ListCourses())}
func(h *progressionHandler)enroll(w http.ResponseWriter,r *http.Request){var q struct{CharacterID string `json:"characterId"`;CourseID string `json:"courseId"`};if json.NewDecoder(r.Body).Decode(&q)!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return};if err:=h.s.Enroll(q.CharacterID,q.CourseID);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,200,h.s.Get(q.CharacterID))}
