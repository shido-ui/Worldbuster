package api

import("encoding/json";"net/http";"github.com/shido-ui/Worldbuster/backend/internal/inventory")

type inventoryHandler struct{service *inventory.Service}
func newInventoryHandler(s *inventory.Service)*inventoryHandler{return &inventoryHandler{service:s}}
func(h *inventoryHandler)get(w http.ResponseWriter,r *http.Request){id:=r.URL.Query().Get("characterId");if id==""{writeJSON(w,400,map[string]string{"error":"characterId is required"});return};writeJSON(w,200,h.service.Get(id))}
func(h *inventoryHandler)add(w http.ResponseWriter,r *http.Request){var req struct{CharacterID,ItemID string;Quantity int};if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{writeJSON(w,400,map[string]string{"error":"invalid json"});return};if err:=h.service.Add(req.CharacterID,req.ItemID,req.Quantity);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return};writeJSON(w,200,h.service.Get(req.CharacterID))}
