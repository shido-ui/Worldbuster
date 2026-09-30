package api

import (
 "context"
 "time"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
)

type AuthBackend interface {
 Register(context.Context,string,string)(auth.PublicAccount,error)
 Authenticate(context.Context,string,string)(auth.PublicAccount,error)
 CreateSession(context.Context,string,time.Duration)(auth.Session,error)
 ResolveSession(context.Context,string)(auth.Session,error)
 RevokeSession(context.Context,string) error
}

type MemoryAuthBackend struct{ Service *auth.Service }

func (b MemoryAuthBackend) Register(_ context.Context,u,p string)(auth.PublicAccount,error){return b.Service.Register(u,p)}
func (b MemoryAuthBackend) Authenticate(_ context.Context,u,p string)(auth.PublicAccount,error){return b.Service.Authenticate(u,p)}
func (b MemoryAuthBackend) CreateSession(_ context.Context,id string,ttl time.Duration)(auth.Session,error){return b.Service.CreateSession(id,ttl)}
func (b MemoryAuthBackend) ResolveSession(_ context.Context,id string)(auth.Session,error){return b.Service.ResolveSession(id)}
func (b MemoryAuthBackend) RevokeSession(_ context.Context,id string)error{b.Service.RevokeSession(id);return nil}
