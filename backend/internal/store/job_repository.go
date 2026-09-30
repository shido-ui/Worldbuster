package store

import (
 "context"
 "database/sql"
)

type JobRepository struct{DB *DB}

type JobRecord struct{ID string; Name string; Department string; BaseSalary int64; RequiredLevel int; RequiredStat int}

func(r JobRepository) List(ctx context.Context)([]JobRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,name,department,base_salary,required_level,required_stat FROM jobs ORDER BY name");if err!=nil{return nil,err};defer rows.Close()
 out:=[]JobRecord{}
 for rows.Next(){var j JobRecord;if err:=rows.Scan(&j.ID,&j.Name,&j.Department,&j.BaseSalary,&j.RequiredLevel,&j.RequiredStat);err!=nil{return nil,err};out=append(out,j)}
 return out,rows.Err()
}

func(r JobRepository) Employment(ctx context.Context,playerID string)(sql.NullString,error){
 var id sql.NullString
 err:=r.DB.SQL.QueryRowContext(ctx,"SELECT job_id::text FROM player_employment WHERE player_id=$1",playerID).Scan(&id)
 return id,err
}
