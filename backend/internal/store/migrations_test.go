package store

import (
	"context"
	"testing"
)

func TestMigrationRunnerSymbol(t *testing.T) {
	var _ func(context.Context, *DB, string) error = ApplyMigrations
}
