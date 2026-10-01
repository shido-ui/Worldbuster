package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type ReputationAPI struct {
	Repo    store.SocialRepository
	Context *PlayerContext
}

func (a *ReputationAPI) Get(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	x, err := a.Repo.Reputation(r.Context(), p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "reputation unavailable"})
		return
	}
	writeJSON(w, 200, x)
}
func (a *ReputationAPI) Change(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var c store.ReputationChange
	if json.NewDecoder(r.Body).Decode(&c) != nil || c.Source == "" || c.Reason == "" {
		writeJSON(w, 400, map[string]string{"error": "source and reason required"})
		return
	}
	if c.PublicDelta > 100 || c.PublicDelta < -100 || c.TrustDelta > 100 || c.TrustDelta < -100 || c.NotorietyDelta > 100 || c.NotorietyDelta < -100 {
		writeJSON(w, 400, map[string]string{"error": "reputation change too large"})
		return
	}
	x, err := a.Repo.ChangeReputation(r.Context(), p.ID, c)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "reputation change rejected"})
		return
	}
	writeJSON(w, 200, x)
}
