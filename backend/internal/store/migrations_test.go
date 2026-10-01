package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestMigrationRunnerSymbol(t *testing.T) {
	var _ = ApplyMigrations
}

func TestApplyMigrationsFreshPostgres(t *testing.T) {
	dsn := os.Getenv("WORLDBUSTER_DATABASE_URL")
	if dsn == "" {
		t.Skip("WORLDBUSTER_DATABASE_URL is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	db := New(sqlDB)
	migrationDir := os.Getenv("WORLDBUSTER_MIGRATIONS_DIR")
	if migrationDir == "" {
		migrationDir = filepath.Join("database", "migrations")
	}
	if info, statErr := os.Stat(migrationDir); statErr != nil || !info.IsDir() {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		migrationDir = ""
		for i := 0; i < 8; i++ {
			candidate := filepath.Join(wd, "database", "migrations")
			if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
				migrationDir = candidate
				break
			}
			wd = filepath.Dir(wd)
		}
		if migrationDir == "" {
			t.Fatal("could not locate database/migrations")
		}
	}

	if err := ApplyMigrations(ctx, db, migrationDir); err != nil {
		t.Fatalf("apply migrations to fresh database: %v", err)
	}

	expected := 0
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			raw, readErr := os.ReadFile(filepath.Join(migrationDir, entry.Name()))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if strings.TrimSpace(string(raw)) != "" {
				expected++
			}
		}
	}

	var applied int
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&applied); err != nil {
		t.Fatal(err)
	}
	if applied != expected {
		t.Fatalf("migration count mismatch: applied=%d expected=%d", applied, expected)
	}

	if err := ApplyMigrations(ctx, db, migrationDir); err != nil {
		t.Fatalf("re-apply migrations: %v", err)
	}

	var appliedAgain int
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&appliedAgain); err != nil {
		t.Fatal(err)
	}
	if appliedAgain != applied {
		t.Fatalf("migration runner is not idempotent: first=%d second=%d", applied, appliedAgain)
	}
}
