package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyMigrations(ctx context.Context, db *DB, dir string) error {
	if db == nil || db.SQL == nil {
		return fmt.Errorf("database is not configured")
	}
	if dir == "" {
		dir = "database/migrations"
	}
	if _, err := db.SQL.ExecContext(ctx, `SELECT pg_advisory_lock(hashtext('worldbuster:migrations'))`); err != nil {
		return err
	}
	defer db.SQL.ExecContext(context.Background(), `SELECT pg_advisory_unlock(hashtext('worldbuster:migrations'))`)

	if _, err := db.SQL.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			checksum TEXT,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}
	if _, err := db.SQL.ExecContext(ctx, `
		ALTER TABLE schema_migrations
		ADD COLUMN IF NOT EXISTS checksum TEXT
	`); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		sqlText := strings.TrimSpace(string(raw))
		if sqlText == "" {
			continue
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(sqlText)))

		var storedChecksum sql.NullString
		err = db.SQL.QueryRowContext(
			ctx,
			"SELECT checksum FROM schema_migrations WHERE version=$1",
			name,
		).Scan(&storedChecksum)
		switch err {
		case nil:
			if storedChecksum.Valid && storedChecksum.String != "" {
				if storedChecksum.String != checksum {
					return fmt.Errorf("migration %s checksum mismatch: database=%s files=%s", name, storedChecksum.String, checksum)
				}
				continue
			}
			if _, err := db.SQL.ExecContext(
				ctx,
				"UPDATE schema_migrations SET checksum=$2 WHERE version=$1",
				name, checksum,
			); err != nil {
				return fmt.Errorf("backfill migration %s checksum: %w", name, err)
			}
			continue
		case sql.ErrNoRows:
			// Apply below.
		default:
			return err
		}

		if err := db.WithTx(ctx, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, sqlText); err != nil {
				return err
			}
			_, err := tx.ExecContext(
				ctx,
				"INSERT INTO schema_migrations(version, checksum) VALUES($1,$2)",
				name, checksum,
			)
			return err
		}); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}

	return nil
}
