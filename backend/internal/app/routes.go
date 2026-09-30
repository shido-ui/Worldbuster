package app

import (
	"context"

	"github.com/shido-ui/Worldbuster/backend/internal/api"
	"github.com/shido-ui/Worldbuster/backend/internal/store"
)

type Dependencies struct {
	DB          *store.DB
	Auth        api.AuthBackend
	World       interface{ Snapshot() any }
	Travel      any
	Inventory   any
	Economy     any
	Organization any
	Job         any
	Progression any
	Social      any
	Events      any
}

func ConfigureRouter(
	db *store.DB,
	authBackend api.AuthBackend,
	worldService interface{ Snapshot() any },
	travelService interface{},
	inventoryService interface{},
	economyService interface{},
	organizationService interface{},
	jobService interface{},
	progressionService interface{},
	socialService interface{},
	eventsService interface{},
) *api.Router {
	_ = context.Background()
	_ = worldService
	_ = travelService
	_ = inventoryService
	_ = economyService
	_ = organizationService
	_ = jobService
	_ = progressionService
	_ = socialService
	_ = eventsService
	_ = db
	_ = authBackend
	return nil
}
