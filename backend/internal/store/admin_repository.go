package store

import (
	"context"
	"database/sql"
	"encoding/json"
)

type AdminRepository struct{ DB *DB }

func (r AdminRepository) IsAdmin(ctx context.Context, accountID string) bool {
	if r.DB == nil || r.DB.SQL == nil {
		return false
	}
	var ok bool
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM admin_roles WHERE account_id=$1)", accountID).Scan(&ok)
	return err == nil && ok
}
func (r AdminRepository) SetFlag(ctx context.Context, accountID, key string, value any) error {
	if !r.IsAdmin(ctx, accountID) {
		return sql.ErrNoRows
	}
	b, _ := json.Marshal(value)
	return r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, "INSERT INTO world_control_flags(key,value,updated_by) VALUES($1,$2,$3) ON CONFLICT(key) DO UPDATE SET value=EXCLUDED.value,updated_by=EXCLUDED.updated_by,updated_at=NOW()", key, b, accountID); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, "INSERT INTO world_control_audit(account_id,action,target_type,target_id,payload) VALUES($1,$2,$3,$4,$5)", accountID, "SET_FLAG", "WORLD_FLAG", key, b)
		return err
	})
}
func (r AdminRepository) Flags(ctx context.Context) (map[string]json.RawMessage, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT key,value FROM world_control_flags ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := map[string]json.RawMessage{}
	for rows.Next() {
		var k string
		var b []byte
		if err := rows.Scan(&k, &b); err != nil {
			return nil, err
		}
		out[k] = json.RawMessage(b)
	}
	return out, rows.Err()
}
