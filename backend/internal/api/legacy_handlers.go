package api

import (
 "encoding/json"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
 "github.com/shido-ui/Worldbuster/backend/internal/inventory"
 "github.com/shido-ui/Worldbuster/backend/internal/progression"
)

type inventoryHandler struct{s *inventory.Service}
func newInventoryHandler(s *inventory.Service)*inventoryHandler{return &inventoryHandler{s:s}}
func(h *inventoryHandler)get(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");if id==""{writeJSON(w,400,map[string]string{"error":"characterId required"});return};writeJSON(w,200,h.s.Get(id))}
func(h *inventoryHandler)add(w http.ResponseWriter,r *http.Request){var q struct{CharacterID string `json:"characterId"`;ItemID string `json:"itemId"`;Quantity int `json:"quantity"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.CharacterID==""||q.ItemID==""||q.Quantity<=0{writeJSON(w,400,map[string]string{"error":"invalid request"});return};if err:=h.s.Add(q.CharacterID,q.ItemID,q.Quantity);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,200,h.s.Get(q.CharacterID))}

type economyHandler struct{s *economy.Service}
func newEconomyHandler(s *economy.Service)*economyHandler{return &economyHandler{s:s}}
func(h *economyHandler)balance(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("accountId");if id==""{writeJSON(w,400,map[string]string{"error":"accountId required"});return};a,err:=h.s.Balance(id);if err!=nil{writeJSON(w,404,map[string]string{"error":err.Error()});return};writeJSON(w,200,a)}
func(h *economyHandler)credit(w http.ResponseWriter,r *http.Request){var q struct{AccountID string `json:"accountId"`;Amount int64 `json:"amount"`;Reason string `json:"reason"`;Reference string `json:"referenceId"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.AccountID==""||q.Amount<=0{writeJSON(w,400,map[string]string{"error":"invalid request"});return};e,err:=h.s.Credit(q.AccountID,q.Amount,q.Reason,q.Reference);if err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,200,e)}
func(h *economyHandler)history(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("accountId");if id==""{writeJSON(w,400,map[string]string{"error":"accountId required"});return};writeJSON(w,200,h.s.History(id))}

type progressionHandler struct{s *progression.Service}
func newProgressionHandler(s *progression.Service)*progressionHandler{return &progressionHandler{s:s}}
func(h *progressionHandler)get(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");if id==""{writeJSON(w,400,map[string]string{"error":"characterId required"});return};writeJSON(w,200,h.s.Get(id))}
func(h *progressionHandler)courses(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,h.s.ListCourses())}
func(h *progressionHandler)enroll(w http.ResponseWriter,r *http.Request){var q struct{CharacterID string `json:"characterId"`;CourseID string `json:"courseId"`};if json.NewDecoder(r.Body).Decode(&q)!=nil||q.CharacterID==""||q.CourseID==""{writeJSON(w,400,map[string]string{"error":"invalid request"});return};if err:=h.s.Enroll(q.CharacterID,q.CourseID);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,200,map[string]string{"status":"enrolled"})}
