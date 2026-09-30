package store

import (
 "context"
 "database/sql"
 "errors"
 "strings"
 "time"

 "github.com/lib/pq"
 "golang.org/x/crypto/bcrypt"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
 "github.com/shido-ui/Worldbuster/backend/internal/player"
)

type DatabaseAuthService struct {
 DB *DB
 Accounts AccountRepository
 Sessions SessionRepository
 Players PlayerRepository
 Economy EconomyRepository
}

func NewDatabaseAuthService(db *DB)*DatabaseAuthService{
 return &DatabaseAuthService{DB:db,Accounts:AccountRepository{DB:db},Sessions:SessionRepository{DB:db},Players:PlayerRepository{DB:db},Economy:EconomyRepository{DB:db}}
}

func(s *DatabaseAuthService) Register(ctx context.Context,username,password string)(auth.PublicAccount,error){
 username=strings.TrimSpace(username)
 if len(username)<3||len(username)>24{return auth.PublicAccount{},auth.ErrInvalidUsername}
 if len(password)<8{return auth.PublicAccount{},auth.ErrInvalidPassword}
 hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if err!=nil{return auth.PublicAccount{},err}
 var a auth.Account
 err=s.DB.WithTx(ctx,func(tx *sql.Tx)error{
  if err:=tx.QueryRowContext(ctx,"INSERT INTO accounts(username,password_hash) VALUES($1,$2) RETURNING id::text,username,password_hash,created_at,updated_at",username,string(hash)).Scan(&a.ID,&a.Username,&a.PasswordHash,&a.CreatedAt,&a.UpdatedAt);err!=nil{return err}
  profile:=player.Profile{ID:a.ID,AccountID:a.ID,DisplayName:a.Username,Level:1,XP:0,Cash:0,Energy:100,Strength:1,Defense:1,Speed:1,Intelligence:1,Endurance:1}
  if err:=tx.QueryRowContext(ctx,"INSERT INTO player_profiles(id,account_id,display_name,level,xp,cash,energy,strength,defense,speed,intelligence,endurance) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id",a.ID,a.ID,profile.DisplayName,profile.Level,profile.XP,profile.Cash,profile.Energy,profile.Strength,profile.Defense,profile.Speed,profile.Intelligence,profile.Endurance).Scan(&profile.ID);err!=nil{return err}
  _,err:=tx.ExecContext(ctx,"INSERT INTO economy_accounts(owner_account_id,currency,balance) VALUES($1,'WBX',0)",a.ID)
  return err
 })
 if err!=nil{
  var pe *pq.Error
  if errors.As(err,&pe)&&pe.Code=="23505"{return auth.PublicAccount{},auth.ErrUsernameTaken}
  return auth.PublicAccount{},err
 }
 return a.Public(),nil
}

func(s *DatabaseAuthService) Authenticate(ctx context.Context,username,password string)(auth.PublicAccount,error){
 a,err:=s.Accounts.ByUsername(ctx,strings.TrimSpace(username))
 if err!=nil||bcrypt.CompareHashAndPassword([]byte(a.PasswordHash),[]byte(password))!=nil{return auth.PublicAccount{},auth.ErrInvalidCredentials}
 return a.Public(),nil
}

func(s *DatabaseAuthService) CreateSession(ctx context.Context,accountID string,ttl time.Duration)(auth.Session,error){
 if ttl<=0{ttl=24*time.Hour}
 if _,err:=s.Accounts.ByID(ctx,accountID);err!=nil{return auth.Session{},auth.ErrInvalidCredentials}
 return s.Sessions.Create(ctx,accountID,time.Now().UTC().Add(ttl))
}

func(s *DatabaseAuthService) ResolveSession(ctx context.Context,id string)(auth.Session,error){
 if strings.TrimSpace(id)==""{return auth.Session{},auth.ErrSessionNotFound}
 v,err:=s.Sessions.Get(ctx,id)
 if err!=nil{return auth.Session{},auth.ErrSessionNotFound}
 return v,nil
}

func(s *DatabaseAuthService) RevokeSession(ctx context.Context,id string)error{
 if strings.TrimSpace(id)==""{return errors.New("invalid session")}
 return s.Sessions.Revoke(ctx,id)
}
