package store

import (
 "context"
 "database/sql"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
)

type EconomyRepository struct{DB *DB}

func(r EconomyRepository) GetByAccountID(ctx context.Context,accountID string)(economy.Account,error){
 var a economy.Account
 err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,balance,currency FROM economy_accounts WHERE owner_account_id=$1",accountID).Scan(&a.ID,&a.Balance,&a.Currency)
 if err==sql.ErrNoRows{return economy.Account{},ErrNotFound};return a,err
}

func(r EconomyRepository) CreateForAccount(ctx context.Context,accountID string)(economy.Account,error){
 var a economy.Account
 err:=r.DB.SQL.QueryRowContext(ctx,"INSERT INTO economy_accounts(owner_account_id,currency,balance) VALUES($1,'WBX',0) ON CONFLICT(owner_account_id) DO UPDATE SET owner_account_id=EXCLUDED.owner_account_id RETURNING id::text,balance,currency",accountID).Scan(&a.ID,&a.Balance,&a.Currency)
 return a,err
}
