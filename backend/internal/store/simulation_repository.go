package store

import("context";"time")

type SimulationMemoryRecord struct{ID string;CharacterID string;MemoryType string;TargetID string;Event string;Importance int;Sentiment int;CreatedAt time.Time}
type SimulationActionRecord struct{ID string;CharacterID string;ActionType string;TargetID string;DecisionScore int;Executed bool;CreatedAt time.Time}
type SimulationRepository struct{DB *DB}
func(r SimulationRepository) RecordMemory(ctx context.Context,characterID,memoryType,targetID,event string,importance,sentiment int)error{
 _,err:=r.DB.SQL.ExecContext(ctx,"INSERT INTO simulated_character_memory(character_id,memory_type,target_id,event,importance,sentiment) VALUES($1,$2,$3,$4,$5,$6)",characterID,memoryType,targetID,event,importance,sentiment);return err
}
func(r SimulationRepository) RecordAction(ctx context.Context,characterID,actionType,targetID string,score int,executed bool)error{
 _,err:=r.DB.SQL.ExecContext(ctx,"INSERT INTO simulated_character_actions(character_id,action_type,target_id,decision_score,executed) VALUES($1,$2,$3,$4,$5)",characterID,actionType,targetID,score,executed);return err
}
func(r SimulationRepository) Memories(ctx context.Context,characterID string,limit int)([]SimulationMemoryRecord,error){
 if limit<=0||limit>200{limit=50};rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,character_id::text,memory_type,COALESCE(target_id,''),event,importance,sentiment,created_at FROM simulated_character_memory WHERE character_id=$1 ORDER BY created_at DESC LIMIT $2",characterID,limit);if err!=nil{return nil,err};defer rows.Close();out:=[]SimulationMemoryRecord{};for rows.Next(){var x SimulationMemoryRecord;if err:=rows.Scan(&x.ID,&x.CharacterID,&x.MemoryType,&x.TargetID,&x.Event,&x.Importance,&x.Sentiment,&x.CreatedAt);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()
}
