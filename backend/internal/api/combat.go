package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type CombatAPI struct {
	Repo    store.CombatRepository
	Context *PlayerContext
}

func (a *CombatAPI) Resolve(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var q struct {
		DefenderID string `json:"defenderId"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil || q.DefenderID == "" {
		writeJSON(w, 400, map[string]string{"error": "defenderId required"})
		return
	}
	x, err := a.Repo.Resolve(r.Context(), p.ID, q.DefenderID)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "combat unavailable"})
		return
	}
	writeJSON(w, 200, x)
}
func (a *CombatAPI) History(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	x, err := a.Repo.History(r.Context(), p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "combat history unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"matches": x})
}
