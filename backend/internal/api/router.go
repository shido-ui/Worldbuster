package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

type Router struct {
	world *world.Service
}

func NewRouter(w *world.Service) *Router { return &Router{world: w} }

func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", r.health)
	mux.HandleFunc("/api/v1/world", r.worldState)
	return mux
}

func (r *Router) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":"ok", "service":"worldbuster", "time":time.Now().UTC(),
	})
}

func (r *Router) worldState(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, r.world.Snapshot())
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
