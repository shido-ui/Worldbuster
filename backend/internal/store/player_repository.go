package store

import("context";"database/sql";"errors";"github.com/shido-ui/Worldbuster/backend/internal/player")
type PlayerRepository struct{DB *DB}
func(r PlayerRepository)Create(ctx context.Context,p player.Profile)(player.Profile,error){var o player.Profile;err:=r.DB.SQL.QueryRowContext(ctx,"INSERT INTO player_profiles(id,account_id,display_name,level,xp,cash,energy,strength,defense,speed,intelligence,endurance) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id::text,account_id::text,display_name,level,xp,cash,energy,strength,defense,speed,intelligence,endurance,education,location_id,updated_at",p.ID,p.AccountID,p.DisplayName,p.Level,p.XP,p.Cash,p.Energy,p.Strength,p.Defense,p.Speed,p.Intelligence,p.Endurance).Scan(&o.ID,&o.AccountID,&o.DisplayName,&o.Level,&o.XP,&o.Cash,&o.Energy,&o.Strength,&o.Defense,&o.Speed,&o.Intelligence,&o.Endurance,&o.Education,&o.LocationID,&o.UpdatedAt);return o,err}
func(r PlayerRepository)GetByID(ctx context.Context,id string)(player.Profile,error){var p player.Profile;err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,account_id::text,display_name,level,xp,cash,energy,strength,defense,speed,intelligence,endurance,education,location_id,updated_at FROM player_profiles WHERE id=$1",id).Scan(&p.ID,&p.AccountID,&p.DisplayName,&p.Level,&p.XP,&p.Cash,&p.Energy,&p.Strength,&p.Defense,&p.Speed,&p.Intelligence,&p.Endurance,&p.Education,&p.LocationID,&p.UpdatedAt);if errors.Is(err,sql.ErrNoRows){return player.Profile{},ErrNotFound};return p,err}


func(r PlayerRepository) GetByAccountID(ctx context.Context,accountID string)(player.Profile,error){
 var p player.Profile
 err:=r.DB.SQL.QueryRowContext(ctx,"SELECT id::text,account_id::text,display_name,level,xp,cash,energy,strength,defense,speed,intelligence,endurance,education,location_id,updated_at FROM player_profiles WHERE account_id=$1",accountID).Scan(&p.ID,&p.AccountID,&p.DisplayName,&p.Level,&p.XP,&p.Cash,&p.Energy,&p.Strength,&p.Defense,&p.Speed,&p.Intelligence,&p.Endurance,&p.Education,&p.LocationID,&p.UpdatedAt)
 if errors.Is(err,sql.ErrNoRows){return player.Profile{},ErrNotFound}
 return p,err
}

func(r PlayerRepository)SetLocation(ctx context.Context,playerID,location string)error{_,err:=r.DB.SQL.ExecContext(ctx,"UPDATE player_profiles SET location_id=$2,updated_at=NOW() WHERE id=$1",playerID,location);return err}
