package api

import (
	"github.com/shido-ui/Worldbuster/backend/internal/auth"
	"github.com/shido-ui/Worldbuster/backend/internal/economy"
	"github.com/shido-ui/Worldbuster/backend/internal/events"
	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/organization"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/social"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlayerStateRequiresAuthentication(t *testing.T) {
	backend := MemoryAuthBackend{Service: auth.NewService()}
	r := NewRouter(world.NewService(), backend, world.NewTravelService(), inventory.NewService(), economy.NewService(), organization.NewService(), job.NewService(), progression.NewService(), social.NewService(), events.NewService())
	r.WithPlayerContext(NewPlayerContext(nil, backend))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/player/state", nil)
	rr := httptest.NewRecorder()
	r.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
