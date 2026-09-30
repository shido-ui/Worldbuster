package auth

import (
 "context"
 "errors"
 "strings"
 "time"
 "golang.org/x/crypto/bcrypt"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type PostgresService struct {
 Accounts store.AccountRepository
 Sessions store.SessionRepository
}

func NewPostgresService(accounts store.AccountRepository,sessions store.SessionRepository)*PostgresService{
 return &PostgresService{Accounts:accounts,Sessions:sessions}
}

func(s *PostgresService) Register(ctx context.Context,username,password string)(PublicAccount,error){
 username=strings.TrimSpace(username)
 if len(username)<3||len(username)>24{return PublicAccount{},ErrInvalidUsername}
 if len(password)<8{return PublicAccount{},ErrInvalidPassword}
 hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if err!=nil{return PublicAccount{},err}
 a,err:=s.Accounts.Create(ctx,username,string(hash))
 if err!=nil{return PublicAccount{},err}
 return a.Public(),nil
}

func(s *PostgresService) Authenticate(ctx context.Context,username,password string)(PublicAccount,error){
 a,err:=s.Accounts.ByUsername(ctx,username)
 if err!=nil||bcrypt.CompareHashAndPassword([]byte(a.PasswordHash),[]byte(password))!=nil{return PublicAccount{},ErrInvalidCredentials}
 return a.Public(),nil
}

func(s *PostgresService) CreateSession(ctx context.Context,accountID string,ttl time.Duration)(Session,error){
 if ttl<=0{ttl=24*time.Hour}
 if _,err:=s.Accounts.ByID(ctx,accountID);err!=nil{return Session{},ErrInvalidCredentials}
 return s.Sessions.Create(ctx,accountID,time.Now().UTC().Add(ttl))
}

func(s *PostgresService) ResolveSession(ctx context.Context,id string)(Session,error){
 if strings.TrimSpace(id)==""{return Session{},ErrSessionNotFound}
 session,err:=s.Sessions.Get(ctx,id)
 if err!=nil{return Session{},ErrSessionNotFound}
 return session,nil
}

func(s *PostgresService) RevokeSession(ctx context.Context,id string)error{
 if strings.TrimSpace(id)==""{return errors.New("invalid session")}
 return s.Sessions.Revoke(ctx,id)
}
