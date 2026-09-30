package store

import ("context";"database/sql";"errors";"time")

type ContractRecord struct {ID,IssuerType,Title,Description,TargetType,Status string; IssuerID *string; TargetValue,RewardCash,RewardXP int64; ExpiresAt *time.Time; CreatedAt time.Time}
type ContractClaimRecord struct {ContractID,PlayerID string; AcceptedAt time.Time; CompletedAt *time.Time}
var ErrContractState=errors.New("invalid contract state")

type ContractRepository struct{DB *DB}

func(r ContractRepository) ListOpen(ctx context.Context)([]ContractRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,issuer_type,issuer_id::text,title,description,target_type,target_value,reward_cash,reward_xp,status,expires_at,created_at FROM contracts WHERE status='OPEN' AND (expires_at IS NULL OR expires_at>NOW()) ORDER BY created_at DESC LIMIT 100")
 if err!=nil{return nil,err};defer rows.Close();out:=[]ContractRecord{}
 for rows.Next(){var x ContractRecord;var issuer,expiry sql.NullString;if err:=rows.Scan(&x.ID,&x.IssuerType,&issuer,&x.Title,&x.Description,&x.TargetType,&x.TargetValue,&x.RewardCash,&x.RewardXP,&x.Status,&expiry,&x.CreatedAt);err!=nil{return nil,err};if issuer.Valid{x.IssuerID=&issuer.String};if expiry.Valid{if t,e:=time.Parse(time.RFC3339,expiry.String);e==nil{x.ExpiresAt=&t}};out=append(out,x)};return out,rows.Err()
}

func(r ContractRepository) Accept(ctx context.Context,id,playerID string)(ContractClaimRecord,error){
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return ContractClaimRecord{},err};defer tx.Rollback()
 var x ContractClaimRecord
 err=tx.QueryRowContext(ctx,"INSERT INTO contract_claims(contract_id,player_id) SELECT id,$2 FROM contracts WHERE id=$1 AND status='OPEN' AND (expires_at IS NULL OR expires_at>NOW()) RETURNING contract_id::text,player_id::text,accepted_at,completed_at",id,playerID).Scan(&x.ContractID,&x.PlayerID,&x.AcceptedAt,&x.CompletedAt)
 if err==sql.ErrNoRows{return x,ErrContractState};if err!=nil{return x,err}
 if _,err=tx.ExecContext(ctx,"UPDATE contracts SET status='ACCEPTED' WHERE id=$1 AND status='OPEN'",id);err!=nil{return x,err}
 if err=tx.Commit();err!=nil{return x,err};return x,nil
}
