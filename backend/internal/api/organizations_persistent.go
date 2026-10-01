package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type PersistentOrganizationsAPI struct {
	Repo    store.OrganizationRepository
	Context *PlayerContext
}

func (a *PersistentOrganizationsAPI) List(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.List(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "organizations unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"organizations": x})
}
func (a *PersistentOrganizationsAPI) Join(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var q struct {
		OrganizationID string `json:"organizationId"`
		Role           string `json:"role"`
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil || q.OrganizationID == "" {
		writeJSON(w, 400, map[string]string{"error": "organizationId required"})
		return
	}
	if q.Role == "" {
		q.Role = "member"
	}
	m, err := a.Repo.Join(r.Context(), q.OrganizationID, p.ID, q.Role)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "unable to join organization"})
		return
	}
	writeJSON(w, 201, m)
}
func (a *PersistentOrganizationsAPI) Reputation(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	id := r.URL.Query().Get("organizationId")
	if id == "" {
		writeJSON(w, 400, map[string]string{"error": "organizationId required"})
		return
	}
	n, err := a.Repo.FactionReputation(r.Context(), id, p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "faction reputation unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"organizationId": id, "playerId": p.ID, "score": n})
}
