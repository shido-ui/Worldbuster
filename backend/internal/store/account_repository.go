package store

import("context";"database/sql";"errors";"time";"github.com/shido-ui/Worldbuster/backend/internal/auth")
var ErrNotFound=errors.New("not found")
type AccountRepository struct{DB *DB}
func(r AccountRepository)Create(ctx context.Context,username,passwordHash string)(auth.Account,error){var a auth.Account;err:=r.DB.SQL.QueryRowContext(ctx,"INSERT INTO accounts(username,password_hash) VALUES($1,$2) RETURNING id::text,username,password_hash,created_at,updated_at",username,passwordHash).Scan(&a.ID,&a.Username,&a.PasswordHash,&a.CreatedAt,&a.UpdatedAt);return a,err}
func(r AccountRepository)ByUsername(ctx context.Context,username string)(auth.Account,error){var a auth.Account;err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,username,password_hash,created_at,updated_at FROM accounts WHERE lower(username)=lower($1)",username).Scan(&a.ID,&a.Username,&a.PasswordHash,&a.CreatedAt,&a.UpdatedAt);if errors.Is(err,sql.ErrNoRows){return auth.Account{},ErrNotFound};return a,err}
func(r AccountRepository)ByID(ctx context.Context,id string)(auth.Account,error){var a auth.Account;err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,username,password_hash,created_at,updated_at FROM accounts WHERE id=$1",id).Scan(&a.ID,&a.Username,&a.PasswordHash,&a.CreatedAt,&a.UpdatedAt);if errors.Is(err,sql.ErrNoRows){return auth.Account{},ErrNotFound};return a,err}
type SessionRepository struct{DB *DB}
func(r SessionRepository)Create(ctx context.Context,accountID string,expires time.Time)(auth.Session,error){var s auth.Session;err:=r.DB.SQL.QueryRowContext(ctx,"INSERT INTO sessions(account_id,expires_at) VALUES($1,$2) RETURNING id::text,account_id::text,expires_at,created_at",accountID,expires).Scan(&s.ID,&s.AccountID,&s.ExpiresAt,&s.CreatedAt);return s,err}
func(r SessionRepository)Get(ctx context.Context,id string)(auth.Session,error){var s auth.Session;err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,account_id::text,expires_at,created_at FROM sessions WHERE id=$1 AND expires_at>NOW()",id).Scan(&s.ID,&s.AccountID,&s.ExpiresAt,&s.CreatedAt);if errors.Is(err,sql.ErrNoRows){return auth.Session{},ErrNotFound};return s,err}
func(r SessionRepository)Revoke(ctx context.Context,id string)error{_,err:=r.DB.SQL.ExecContext(ctx,"DELETE FROM sessions WHERE id=$1",id);return err}

func(r AccountRepository) DeleteByID(ctx context.Context,id string) error { _,err:=r.DB.SQL.ExecContext(ctx,"DELETE FROM accounts WHERE id=$1",id); return err }
