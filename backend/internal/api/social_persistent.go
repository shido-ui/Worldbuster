package api

import (
	"database/sql"
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type PersistentSocialAPI struct {
	Repo    store.SocialRepository
	Context *PlayerContext
}

func (a *PersistentSocialAPI) Inbox(w http.ResponseWriter, r *http.Request) {
	if a.Context == nil {
		writeJSON(w, 503, map[string]string{"error": "player persistence unavailable"})
		return
	}
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	items, err := a.Repo.Inbox(r.Context(), p.ID)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "inbox unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"messages": items})
}
func (a *PersistentSocialAPI) Send(w http.ResponseWriter, r *http.Request) {
	if a.Context == nil {
		writeJSON(w, 503, map[string]string{"error": "player persistence unavailable"})
		return
	}
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var req map[string]string
	if json.NewDecoder(r.Body).Decode(&req) != nil || req["toId"] == "" || req["body"] == "" {
		writeJSON(w, 400, map[string]string{"error": "toId and body required"})
		return
	}
	m, err := a.Repo.Send(r.Context(), p.ID, req["toId"], req["body"])
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "message could not be sent"})
		return
	}
	writeJSON(w, 201, m)
}
func (a *PersistentSocialAPI) Relationship(w http.ResponseWriter, r *http.Request) {
	if a.Context == nil {
		writeJSON(w, 503, map[string]string{"error": "player persistence unavailable"})
		return
	}
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	target := r.URL.Query().Get("targetId")
	if target == "" {
		writeJSON(w, 400, map[string]string{"error": "targetId required"})
		return
	}
	rel, err := a.Repo.Relationship(r.Context(), p.ID, target)
	if err == sql.ErrNoRows {
		writeJSON(w, 200, map[string]any{"relationship": nil})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "relationship unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"relationship": rel})
}
func (a *PersistentSocialAPI) SetRelationship(w http.ResponseWriter, r *http.Request) {
	if a.Context == nil {
		writeJSON(w, 503, map[string]string{"error": "player persistence unavailable"})
		return
	}
	_, p, err := a.Context.Resolve(r.Context(), r)
	if err != nil {
		writeJSON(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var req struct {
		TargetID string `json:"targetId"`
		Kind     string `json:"kind"`
		Score    int    `json:"score"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.TargetID == "" || req.Kind == "" {
		writeJSON(w, 400, map[string]string{"error": "targetId and kind required"})
		return
	}
	rel, err := a.Repo.SetRelationship(r.Context(), p.ID, req.TargetID, req.Kind, req.Score)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "relationship could not be updated"})
		return
	}
	writeJSON(w, 200, rel)
}
