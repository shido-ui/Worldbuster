package api

import (
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type AchievementsAPI struct {
	Repo    store.AchievementRepository
	Context *PlayerContext
}

func (a *AchievementsAPI) List(w http.ResponseWriter, r *http.Request) {
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	x, err := a.Repo.List(r.Context(), p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "achievements unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"achievements": x})
}
func (a *AchievementsAPI) Rankings(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = "level"
	}
	x, err := a.Repo.Rankings(r.Context(), id)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "rankings unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"ranking": id, "entries": x})
}
