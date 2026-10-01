package auth

import (
	"context"
	"time"
)

type Provider interface {
	Register(context.Context, string, string) (PublicAccount, error)
	Authenticate(context.Context, string, string) (PublicAccount, error)
	CreateSession(context.Context, string, time.Duration) (Session, error)
	ResolveSession(context.Context, string) (Session, error)
	RevokeSession(context.Context, string) error
}
