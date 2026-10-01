package api

import (
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type AdminAPI struct {
	Repo store.AdminRepository
	Auth AuthBackend
}

func (a *AdminAPI) Flags(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.Flags(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "world controls unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"flags": x})
}
func (a *AdminAPI) SetFlag(w http.ResponseWriter, r *http.Request) {
	session, ok := sessionFromRequest(a.Auth, r)
	if !ok || !a.Repo.IsAdmin(r.Context(), session.AccountID) {
		writeJSON(w, 403, map[string]string{"error": "world operator access required"})
		return
	}
	var req struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Key == "" || len(req.Value) == 0 {
		writeJSON(w, 400, map[string]string{"error": "invalid flag"})
		return
	}
	var v any
	if json.Unmarshal(req.Value, &v) != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid value"})
		return
	}
	if err := a.Repo.SetFlag(r.Context(), session.AccountID, req.Key, v); err != nil {
		writeJSON(w, 500, map[string]string{"error": "control update failed"})
		return
	}
	writeJSON(w, 200, map[string]any{"updated": true, "key": req.Key})
}
