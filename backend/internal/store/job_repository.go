package store

import (
 "context"
 "database/sql"
 "errors"
)

var ErrJobNotFound = errors.New("job not found")
var ErrAlreadyEmployed = errors.New("already employed")
var ErrJobRequirement = errors.New("job requirements not met")

type JobRepository struct{DB *DB}
type JobRecord struct { ID string; Name string; Department string; BaseSalary int64; RequiredLevel int; RequiredStat int }

func(r JobRepository) List(ctx context.Context)([]JobRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,name,department,base_salary,required_level,required_stat FROM jobs ORDER BY required_level, name");if err!=nil{return nil,err};defer rows.Close()
 out:=[]JobRecord{};for rows.Next(){var j JobRecord;if err:=rows.Scan(&j.ID,&j.Name,&j.Department,&j.BaseSalary,&j.RequiredLevel,&j.RequiredStat);err!=nil{return nil,err};out=append(out,j)};return out,rows.Err()
}
func(r JobRepository) Employment(ctx context.Context,playerID string)(sql.NullString,error){var id sql.NullString;err:=r.DB.SQL.QueryRowContext(ctx,"SELECT job_id::text FROM player_employment WHERE player_id=$1",playerID).Scan(&id);return id,err}
func(r JobRepository) Employ(ctx context.Context,playerID,jobID string,level,stat int)(JobRecord,error){
 var j JobRecord
 err:=r.DB.WithTx(ctx,func(tx *sql.Tx)error{
  var existing string
  err:=tx.QueryRowContext(ctx,"SELECT job_id::text FROM player_employment WHERE player_id=$1 FOR UPDATE",playerID).Scan(&existing)
  if err==nil{return ErrAlreadyEmployed};if err!=sql.ErrNoRows{return err}
  err=tx.QueryRowContext(ctx,"SELECT id::text,name,department,base_salary,required_level,required_stat FROM jobs WHERE id=$1",jobID).Scan(&j.ID,&j.Name,&j.Department,&j.BaseSalary,&j.RequiredLevel,&j.RequiredStat)
  if err==sql.ErrNoRows{return ErrJobNotFound};if err!=nil{return err};if level<j.RequiredLevel||stat<j.RequiredStat{return ErrJobRequirement}
  _,err=tx.ExecContext(ctx,"INSERT INTO player_employment(player_id,job_id) VALUES($1,$2)",playerID,jobID);return err
 });return j,err
}