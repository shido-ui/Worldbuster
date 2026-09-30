package store

import (
 "context"
 "errors"
 "strings"
 "time"

 "github.com/lib/pq"
 "golang.org/x/crypto/bcrypt"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
)

type DatabaseAuthService struct {
 Accounts AccountRepository
 Sessions SessionRepository
}

func NewDatabaseAuthService(db *DB)*DatabaseAuthService{
 return &DatabaseAuthService{Accounts:AccountRepository{DB:db},Sessions:SessionRepository{DB:db}}
}

func(s *DatabaseAuthService) Register(ctx context.Context,username,password string)(auth.PublicAccount,error){
 username=strings.TrimSpace(username)
 if len(username)<3||len(username)>24{return auth.PublicAccount{},auth.ErrInvalidUsername}
 if len(password)<8{return auth.PublicAccount{},auth.ErrInvalidPassword}
 hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if err!=nil{return auth.PublicAccount{},err}
 a,err:=s.Accounts.Create(ctx,username,string(hash))
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
