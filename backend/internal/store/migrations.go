package store

import (
 "context"
 "database/sql"
 "fmt"
 "os"
 "path/filepath"
 "sort"
 "strings"
)

func ApplyMigrations(ctx context.Context,db *DB,dir string) error {
 if db==nil||db.SQL==nil{return fmt.Errorf("database is not configured")}
 if dir==""{dir="database/migrations"}
 if _,err:=db.SQL.ExecContext(ctx,`SELECT pg_advisory_lock(hashtext('worldbuster:migrations'))`);err!=nil{return err}
 defer db.SQL.ExecContext(context.Background(),`SELECT pg_advisory_unlock(hashtext('worldbuster:migrations'))`)
 if _,err:=db.SQL.ExecContext(ctx,`CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
 )`);err!=nil{return err}
 entries,err:=os.ReadDir(dir);if err!=nil{return err}
 files:=make([]string,0)
 for _,e:=range entries{
  if !e.IsDir()&&strings.HasSuffix(e.Name(),".sql"){files=append(files,e.Name())}
 }
 sort.Strings(files)
 for _,name:=range files{
  var applied bool
  if err:=db.SQL.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)",name).Scan(&applied);err!=nil{return err}
  if applied{continue}
  raw,err:=os.ReadFile(filepath.Join(dir,name));if err!=nil{return err}
  sqlText:=strings.TrimSpace(string(raw));if sqlText==""{continue}
  if err:=db.WithTx(ctx,func(tx *sql.Tx) error{
   if _,err:=tx.ExecContext(ctx,sqlText);err!=nil{return err}
   _,err:=tx.ExecContext(ctx,"INSERT INTO schema_migrations(version) VALUES($1)",name)
   return err
  });err!=nil{return fmt.Errorf("migration %s: %w",name,err)}
 }
 return nil
}
