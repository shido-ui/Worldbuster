package api

import (
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"net/http"
)

type WorldEventsAPI struct{ Repo store.WorldEventRepository }

func (a *WorldEventsAPI) Definitions(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.Definitions(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "world event definitions unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"events": x})
}
func (a *WorldEventsAPI) Recent(w http.ResponseWriter, r *http.Request) {
	x, err := a.Repo.Recent(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "world events unavailable"})
		return
	}
	writeJSON(w, 200, map[string]any{"events": x})
}
