package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shido-ui/Worldbuster/backend/internal/auth"
	"github.com/shido-ui/Worldbuster/backend/internal/economy"
	"github.com/shido-ui/Worldbuster/backend/internal/events"
	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/organization"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/social"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

func TestTravelEndpointProductionPath(t *testing.T) {
	travel := world.NewTravelService()
	router := NewRouter(world.NewService(), MemoryAuthBackend{Service: auth.NewService()}, travel, inventory.NewService(), economy.NewService(), organization.NewService(), job.NewService(), progression.NewService(), social.NewService(), events.NewService())

	body, err := json.Marshal(map[string]string{"characterId": "travel-e2e", "from": "central", "to": "harbor"})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/world/travel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("travel endpoint returned %d, want %d", rr.Code, http.StatusAccepted)
	}

	var got world.TravelState
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.CharacterID != "travel-e2e" || got.From != "central" || got.To != "harbor" {
		t.Fatalf("unexpected travel state: %+v", got)
	}
	if got.ArrivesAt.IsZero() || !got.ArrivesAt.After(got.StartedAt) {
		t.Fatalf("travel did not receive a valid arrival time: %+v", got)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/world/travel", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("duplicate travel returned %d, want %d", rr.Code, http.StatusBadRequest)
	}
}
