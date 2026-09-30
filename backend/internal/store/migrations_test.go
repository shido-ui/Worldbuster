package store

import ("testing";"context")

func TestMigrationRunnerSymbol(t *testing.T){ var _ func(context.Context,*DB,string) error = ApplyMigrations }
