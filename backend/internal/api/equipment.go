package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type EquipmentAPI struct {
	Repo    store.EquipmentRepository
	Context *PlayerContext
}

func (a *EquipmentAPI) Items(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.ListItems(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "items unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"items": x})
}
func (a *EquipmentAPI) Equipped(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	x, err := a.Repo.ListEquipped(r.Context(), p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "equipment unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"equipment": x})
}
func (a *EquipmentAPI) Equip(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var q struct {
		ItemID string `json:"itemId"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil || q.ItemID == "" {
		writeJSON(w, 400, map[string]string{"error": "itemId required"})
		return
	}
	x, err := a.Repo.Equip(r.Context(), p.ID, q.ItemID, p.Level)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "item cannot be equipped"})
		return
	}
	writeJSON(w, 200, x)
}
