package api

import (
 "context"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
 "github.com/shido-ui/Worldbuster/backend/internal/player"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type PlayerContext struct{ Profile store.PlayerRepository }

func NewPlayerContext(db *store.DB)*PlayerContext{return &PlayerContext{Profile:store.PlayerRepository{DB:db}}}

func(pc *PlayerContext) Resolve(ctx context.Context,r *http.Request)(auth.Session,player.Profile,error){
 session,ok:=sessionFromRequestFromContext(r.Context(),r)
 if !ok{return auth.Session{},player.Profile{},auth.ErrSessionNotFound}
 p,err:=pc.Profile.GetByAccountID(ctx,session.AccountID)
 return session,p,err
}

func sessionFromRequestFromContext(ctx context.Context,r *http.Request)(auth.Session,bool){
 provider,ok:=ctx.Value(authProviderKey{}).(AuthBackend);if !ok{return auth.Session{},false}
 return sessionFromRequest(provider,r)
}
type authProviderKey struct{}
