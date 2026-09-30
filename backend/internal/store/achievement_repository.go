package store

import("context";"database/sql";"time")

type Achievement struct{ID,Title,Description,Category string;Points int;Progress int64;CompletedAt *time.Time}
type RankingEntry struct{PlayerID string;Rank int;Score int64}
type AchievementRepository struct{DB *DB}

func(r AchievementRepository)List(ctx context.Context,playerID string)([]Achievement,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,`SELECT d.id,d.title,d.description,d.category,d.points,COALESCE(p.progress,0),p.completed_at FROM achievement_definitions d LEFT JOIN player_achievements p ON p.achievement_id=d.id AND p.player_id=$1 WHERE d.active ORDER BY d.category,d.id`,playerID);if err!=nil{return nil,err};defer rows.Close()
 out:=[]Achievement{};for rows.Next(){var a Achievement;if err:=rows.Scan(&a.ID,&a.Title,&a.Description,&a.Category,&a.Points,&a.Progress,&a.CompletedAt);err!=nil{return nil,err};out=append(out,a)};return out,rows.Err()
}
func(r AchievementRepository)AddProgress(ctx context.Context,playerID,achievementID string,delta int64)(Achievement,error){
 var a Achievement
 err:=r.DB.WithTx(ctx,func(tx *sql.Tx)error{
  _,err:=tx.ExecContext(ctx,`INSERT INTO player_achievements(player_id,achievement_id,progress) VALUES($1,$2,$3) ON CONFLICT(player_id,achievement_id) DO UPDATE SET progress=player_achievements.progress+$3`,playerID,achievementID,delta);if err!=nil{return err}
  return tx.QueryRowContext(ctx,`UPDATE player_achievements p SET completed_at=CASE WHEN p.completed_at IS NULL AND ((d.id='first-level' AND p.progress>=1) OR (d.id='market-participant' AND p.progress>=1)) THEN NOW() ELSE p.completed_at END FROM achievement_definitions d WHERE p.player_id=$1 AND p.achievement_id=d.id RETURNING d.id,d.title,d.description,d.category,d.points,p.progress,p.completed_at`,playerID).Scan(&a.ID,&a.Title,&a.Description,&a.Category,&a.Points,&a.Progress,&a.CompletedAt)
 });return a,err
}
func(r AchievementRepository)Rankings(ctx context.Context,rankingID string)([]RankingEntry,error){
 metric:="level";if rankingID=="wealth"{metric="cash"};query:=`SELECT id::text,level FROM player_profiles ORDER BY level DESC,xp DESC,id LIMIT 100`;if metric=="cash"{query=`SELECT id::text,cash FROM player_profiles ORDER BY cash DESC,id LIMIT 100`}
 rows,err:=r.DB.SQL.QueryContext(ctx,query);if err!=nil{return nil,err};defer rows.Close();out:=[]RankingEntry{};rank:=1;for rows.Next(){var x RankingEntry;if err:=rows.Scan(&x.PlayerID,&x.Score);err!=nil{return nil,err};x.Rank=rank;rank++;out=append(out,x)};return out,rows.Err()
}
