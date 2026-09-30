package api

import (
 "context"
 "net/http"
 "github.com/shido-ui/Worldbuster/backend/internal/auth"
 "github.com/shido-ui/Worldbuster/backend/internal/player"
 "github.com/shido-ui/Worldbuster/backend/internal/store"
)

type PlayerContext struct { Auth AuthBackend; Profile store.PlayerRepository }
func NewPlayerContext(db *store.DB, authBackend AuthBackend)*PlayerContext{return &PlayerContext{Auth:authBackend,Profile:store.PlayerRepository{DB:db}}}
func(pc *PlayerContext) Resolve(ctx context.Context,r *http.Request)(auth.Session,player.Profile,error){session,ok:=sessionFromRequest(pc.Auth,r);if !ok{return auth.Session{},player.Profile{},auth.ErrSessionNotFound};p,err:=pc.Profile.GetByAccountID(ctx,session.AccountID);return session,p,err}
