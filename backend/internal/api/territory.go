package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type TerritoryAPI struct {
	Repo    store.TerritoryRepository
	Context *PlayerContext
}

func (a *TerritoryAPI) List(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.List(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "territories unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"territories": x})
}
func (a *TerritoryAPI) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, 400, map[string]string{"error": "id required"})
		return
	}
	x, err := a.Repo.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "territory not found"})
		return
	}
	writeJSON(w, 200, x)
}
func (a *TerritoryAPI) AddInfluence(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var q struct {
		TerritoryID    string `json:"territoryId"`
		OrganizationID string `json:"organizationId"`
		Delta          int    `json:"delta"`
		Reason         string `json:"reason"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil || q.TerritoryID == "" || q.OrganizationID == "" || q.Reason == "" {
		writeJSON(w, 400, map[string]string{"error": "territoryId, organizationId and reason required"})
		return
	}
	x, err := a.Repo.AddInfluence(r.Context(), q.TerritoryID, q.OrganizationID, p.ID, q.Delta, q.Reason)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "influence change rejected"})
		return
	}
	writeJSON(w, 200, x)
}
