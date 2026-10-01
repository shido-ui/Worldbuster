package store

import (
	"context"
	"database/sql"

	"github.com/shido-ui/Worldbuster/backend/internal/economy"
	"github.com/shido-ui/Worldbuster/backend/internal/inventory"
	"github.com/shido-ui/Worldbuster/backend/internal/player"
)

type JobStore interface {
	List(context.Context) ([]JobRecord, error)
	Employment(context.Context, string) (sql.NullString, error)
	Employ(context.Context, string, string, int, int, int) (JobRecord, error)
	SettleDueSalaries(context.Context, int) ([]SalaryPayment, error)
}

type PlayerStore interface {
	Create(context.Context, player.Profile) (player.Profile, error)
	GetByID(context.Context, string) (player.Profile, error)
	GetByAccountID(context.Context, string) (player.Profile, error)
	SetLocation(context.Context, string, string) error
}

type EconomyStore interface {
	GetByAccountID(context.Context, string) (economy.Account, error)
	CreateForAccount(context.Context, string) (economy.Account, error)
	CreditByAccountID(context.Context, string, int64, string, string) (economy.LedgerEntry, error)
	LedgerByAccountID(context.Context, string) ([]economy.LedgerEntry, error)
}

type InventoryStore interface {
	Get(context.Context, string) (inventory.Inventory, error)
	Add(context.Context, string, string, int) error
}

var _ JobStore = JobRepository{}
var _ PlayerStore = PlayerRepository{}
var _ EconomyStore = EconomyRepository{}
var _ InventoryStore = InventoryRepository{}
