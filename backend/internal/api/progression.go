package api

import (
	"context"
	"encoding/json"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
	"log"
	"net/http"
)

type PlayerResolver func(ctx context.Context, accountID string) (string, error)

type ProgressionAPI struct {
	Repo          store.ProgressionRepository
	ResolvePlayer PlayerResolver
	Auth          AuthBackend
}

func (a ProgressionAPI) AddXP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if a.Auth == nil {
		http.Error(w, "authentication unavailable", http.StatusServiceUnavailable)
		return
	}
	session, ok := sessionFromRequest(a.Auth, r)
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	playerID := session.AccountID
	if a.ResolvePlayer != nil {
		var err error
		playerID, err = a.ResolvePlayer(r.Context(), session.AccountID)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
	if playerID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	result, err := a.Repo.AddXPAndUnlocks(r.Context(), playerID, req.Amount, nil)
	if err != nil {
		http.Error(w, "progression update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("progression response encoding: %v", err)
	}
}
