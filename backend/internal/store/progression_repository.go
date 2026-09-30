package store

import (
 "context"
 "database/sql"
 "fmt"
 "github.com/shido-ui/Worldbuster/backend/internal/progression"
)

type ProgressionRepository struct{ DB *DB }

type ProgressionCommitResult struct {
 Level int
 XP int64
 Granted []progression.ProgressionReward
}

func(r ProgressionRepository) AddXPAndUnlocks(ctx context.Context, playerID string, amount int64, rewards []progression.ProgressionReward)(ProgressionCommitResult,error){
 if r.DB==nil||r.DB.SQL==nil||playerID==""||amount<=0{return ProgressionCommitResult{},fmt.Errorf("invalid progression commit")}
 var out ProgressionCommitResult
 err:=r.DB.WithTx(ctx,func(tx *sql.Tx) error{
  var level int
  var xp int64
  if err:=tx.QueryRowContext(ctx,`SELECT level,xp FROM player_profiles WHERE id=$1 FOR UPDATE`,playerID).Scan(&level,&xp);err!=nil{return err}
  xp+=amount
  for level<100 {
   need:=int64(level*level*100)
   if xp<need{break}
   xp-=need
   level++
  }
  if _,err:=tx.ExecContext(ctx,`UPDATE player_profiles SET level=$2,xp=$3,updated_at=NOW() WHERE id=$1`,playerID,level,xp);err!=nil{return err}
  for _,reward:=range rewards{
   if reward.ID==""{continue}
   var inserted bool
   err:=tx.QueryRowContext(ctx,`INSERT INTO player_unlocks(player_id,reward_id,reward_type,value,amount)
    VALUES($1,$2,$3,$4,$5) ON CONFLICT(player_id,reward_id) DO NOTHING RETURNING TRUE`,
    playerID,reward.ID,reward.Type,reward.Value,reward.Amount).Scan(&inserted)
   if err==nil&&inserted{out.Granted=append(out.Granted,reward)}
   if err==sql.ErrNoRows{continue}
   if err!=nil{return err}
  }
  out.Level=level;out.XP=xp
  return nil
 })
 return out,err
}
