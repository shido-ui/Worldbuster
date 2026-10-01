package store

import (
	"context"
	"database/sql"
)

type DB struct{ SQL *sql.DB }

func New(db *sql.DB) *DB { return &DB{SQL: db} }

type TxFunc func(*sql.Tx) error

func (d *DB) WithTx(ctx context.Context, fn TxFunc) error {
	tx, err := d.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
