package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shido-ui/Worldbuster/backend/internal/economy"
)

type EconomyRepository struct{ DB *DB }

func (r EconomyRepository) GetByAccountID(ctx context.Context, accountID string) (economy.Account, error) {
	var a economy.Account
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT id::text,balance,currency FROM economy_accounts WHERE owner_account_id=$1", accountID).Scan(&a.ID, &a.Balance, &a.Currency)
	if err == sql.ErrNoRows {
		return economy.Account{}, ErrNotFound
	}
	return a, err
}

func (r EconomyRepository) CreateForAccount(ctx context.Context, accountID string) (economy.Account, error) {
	var a economy.Account
	err := r.DB.SQL.QueryRowContext(ctx, "INSERT INTO economy_accounts(owner_account_id,currency,balance) VALUES($1,'WBX',0) ON CONFLICT(owner_account_id) DO UPDATE SET owner_account_id=EXCLUDED.owner_account_id RETURNING id::text,balance,currency", accountID).Scan(&a.ID, &a.Balance, &a.Currency)
	return a, err
}

func (r EconomyRepository) CreditByAccountID(ctx context.Context, accountID string, amount int64, reason, reference string) (economy.LedgerEntry, error) {
	if amount <= 0 {
		return economy.LedgerEntry{}, errors.New("amount must be positive")
	}
	var e economy.LedgerEntry
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		var id string
		var balance int64
		if err := tx.QueryRowContext(ctx, "SELECT id::text,balance FROM economy_accounts WHERE owner_account_id=$1 FOR UPDATE", accountID).Scan(&id, &balance); err != nil {
			return err
		}
		balance += amount
		if _, err := tx.ExecContext(ctx, "UPDATE economy_accounts SET balance=$1,updated_at=NOW() WHERE id=$2", balance, id); err != nil {
			return err
		}
		return tx.QueryRowContext(ctx, "INSERT INTO economy_ledger(id,account_id,amount,balance_after,reason,reference_id) VALUES(gen_random_uuid(),$1,$2,$3,$4,$5) RETURNING id::text,account_id::text,amount,balance_after,reason,COALESCE(reference_id,''),created_at", id, amount, balance, reason, reference).Scan(&e.ID, &e.AccountID, &e.Amount, &e.BalanceAfter, &e.Reason, &e.ReferenceID, &e.CreatedAt)
	})
	return e, err
}

func (r EconomyRepository) LedgerByAccountID(ctx context.Context, accountID string) ([]economy.LedgerEntry, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT l.id::text,l.account_id::text,l.amount,l.balance_after,l.reason,COALESCE(l.reference_id,''),l.created_at FROM economy_ledger l JOIN economy_accounts a ON a.id=l.account_id WHERE a.owner_account_id=$1 ORDER BY l.created_at DESC LIMIT 100", accountID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []economy.LedgerEntry{}
	for rows.Next() {
		var e economy.LedgerEntry
		if err := rows.Scan(&e.ID, &e.AccountID, &e.Amount, &e.BalanceAfter, &e.Reason, &e.ReferenceID, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
