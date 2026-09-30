package store

import (
 "context"
 "database/sql"
 "errors"
 "time"
 "github.com/shido-ui/Worldbuster/backend/internal/economy"
)

var ErrMissionNotFound=errors.New("mission not found")
var ErrMissionRequirement=errors.New("mission requirement not met")
var ErrMissionState=errors.New("invalid mission state")
var ErrContractNotFound=errors.New("contract not found")

type MissionRecord struct { ID,Code,Title,Description,MissionType,TargetType string; MinLevel int; RewardCash,RewardXP,TargetValue int64; Active bool }
type PlayerMissionRecord struct { ID,MissionID,PlayerID,Status string; Progress int64; AcceptedAt time.Time; CompletedAt *time.Time }

type MissionRepository struct{DB *DB}

func(r MissionRepository) List(ctx context.Context,level int)([]MissionRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,code,title,description,mission_type,min_level,reward_cash,reward_xp,COALESCE(target_type,''),target_value,active FROM missions WHERE active AND min_level<=$1 ORDER BY min_level,title",level)
 if err!=nil{return nil,err};defer rows.Close();out:=[]MissionRecord{}
 for rows.Next(){var m MissionRecord;if err:=rows.Scan(&m.ID,&m.Code,&m.Title,&m.Description,&m.MissionType,&m.MinLevel,&m.RewardCash,&m.RewardXP,&m.TargetType,&m.TargetValue,&m.Active);err!=nil{return nil,err};out=append(out,m)};return out,rows.Err()
}

func(r MissionRepository) Accept(ctx context.Context,missionID,playerID string,level int)(PlayerMissionRecord,error){
 var m MissionRecord
 err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,code,title,description,mission_type,min_level,reward_cash,reward_xp,COALESCE(target_type,''),target_value,active FROM missions WHERE id=$1",missionID).Scan(&m.ID,&m.Code,&m.Title,&m.Description,&m.MissionType,&m.MinLevel,&m.RewardCash,&m.RewardXP,&m.TargetType,&m.TargetValue,&m.Active)
 if err==sql.ErrNoRows{return PlayerMissionRecord{},ErrMissionNotFound};if err!=nil{return PlayerMissionRecord{},err};if !m.Active||level<m.MinLevel{return PlayerMissionRecord{},ErrMissionRequirement}
 var x PlayerMissionRecord
 err=r.DB.SQL.QueryRowContext(ctx,"INSERT INTO player_missions(mission_id,player_id) VALUES($1,$2) RETURNING id::text,mission_id::text,player_id::text,status,progress,accepted_at,completed_at",missionID,playerID).Scan(&x.ID,&x.MissionID,&x.PlayerID,&x.Status,&x.Progress,&x.AcceptedAt,&x.CompletedAt)
 return x,err
}

func(r MissionRepository) Progress(ctx context.Context,missionID,playerID string,amount int64)(PlayerMissionRecord,error){
 if amount<=0{return PlayerMissionRecord{},errors.New("progress must be positive")}
 var x PlayerMissionRecord
 err:=r.DB.SQL.QueryRowContext(ctx,"UPDATE player_missions SET progress=progress+$3 WHERE mission_id=$1 AND player_id=$2 AND status='ACTIVE' RETURNING id::text,mission_id::text,player_id::text,status,progress,accepted_at,completed_at",missionID,playerID,amount).Scan(&x.ID,&x.MissionID,&x.PlayerID,&x.Status,&x.Progress,&x.AcceptedAt,&x.CompletedAt)
 if err==sql.ErrNoRows{return PlayerMissionRecord{},ErrMissionState};return x,err
}

func(r MissionRepository) Complete(ctx context.Context,missionID,playerID string)(PlayerMissionRecord,economy.LedgerEntry,error){
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return PlayerMissionRecord{},economy.LedgerEntry{},err};defer tx.Rollback()
 var progress,target int64;var rewardCash,rewardXP int64;var status string
 err=tx.QueryRowContext(ctx,"SELECT pm.progress,m.target_value,m.reward_cash,m.reward_xp,pm.status FROM player_missions pm JOIN missions m ON m.id=pm.mission_id WHERE pm.mission_id=$1 AND pm.player_id=$2 FOR UPDATE",missionID,playerID).Scan(&progress,&target,&rewardCash,&rewardXP,&status)
 if err==sql.ErrNoRows{return PlayerMissionRecord{},economy.LedgerEntry{},ErrMissionState};if err!=nil{return PlayerMissionRecord{},economy.LedgerEntry{},err};if status!="ACTIVE"||progress<target{return PlayerMissionRecord{},economy.LedgerEntry{},ErrMissionRequirement}
 var e economy.LedgerEntry
 var accountID string
 if err=tx.QueryRowContext(ctx,"SELECT account_id::text FROM player_profiles WHERE id=$1",playerID).Scan(&accountID);err!=nil{return PlayerMissionRecord{},e,err}
 var account,bal int64
 if err=tx.QueryRowContext(ctx,"SELECT id,balance FROM economy_accounts WHERE owner_account_id=$1 FOR UPDATE",accountID).Scan(&account,&bal);err!=nil{return PlayerMissionRecord{},e,err}
 bal+=rewardCash
 if _,err=tx.ExecContext(ctx,"UPDATE economy_accounts SET balance=$1,updated_at=NOW() WHERE id=$2",bal,account);err!=nil{return PlayerMissionRecord{},e,err}
 if err=tx.QueryRowContext(ctx,"INSERT INTO economy_ledger(id,account_id,amount,balance_after,reason,reference_id) VALUES(gen_random_uuid(),$1,$2,$3,'mission.reward',$4) RETURNING id::text,account_id::text,amount,balance_after,reason,COALESCE(reference_id,''),created_at",account,rewardCash,bal,missionID).Scan(&e.ID,&e.AccountID,&e.Amount,&e.BalanceAfter,&e.Reason,&e.ReferenceID,&e.CreatedAt);err!=nil{return PlayerMissionRecord{},e,err}
 now:=time.Now().UTC()
 var x PlayerMissionRecord
 err=tx.QueryRowContext(ctx,"UPDATE player_missions SET status='COMPLETED',completed_at=$1 WHERE mission_id=$2 AND player_id=$3 RETURNING id::text,mission_id::text,player_id::text,status,progress,accepted_at,completed_at",now,missionID,playerID).Scan(&x.ID,&x.MissionID,&x.PlayerID,&x.Status,&x.Progress,&x.AcceptedAt,&x.CompletedAt)
 if err!=nil{return PlayerMissionRecord{},e,err}
 if rewardXP>0 {
  var xp,level int64
  if err=tx.QueryRowContext(ctx,"SELECT xp,level FROM player_profiles WHERE id=$1 FOR UPDATE",playerID).Scan(&xp,&level);err!=nil{return PlayerMissionRecord{},e,err}
  xp+=rewardXP
  for level<100 && xp>=level*level*100 { xp-=level*level*100; level++ }
  if _,err=tx.ExecContext(ctx,"UPDATE player_profiles SET xp=$1,level=$2,updated_at=NOW() WHERE id=$3",xp,level,playerID);err!=nil{return PlayerMissionRecord{},e,err}
 }
 if err=tx.Commit();err!=nil{return PlayerMissionRecord{},e,err}
 return x,e,nil
}
