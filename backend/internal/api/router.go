package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shido-ui/Worldbuster/backend/internal/auth"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

type Router struct { world *world.Service; auth *auth.Service }

func NewRouter(w *world.Service, a *auth.Service) *Router { return &Router{world:w,auth:a} }

func (r *Router) Handler() http.Handler {
	mux:=http.NewServeMux()
	mux.HandleFunc("/health",r.health)
	mux.HandleFunc("/api/v1/world",r.worldState)
	authHandler:=newAuthHandler(r.auth)
	mux.HandleFunc("/api/v1/auth/register",authHandler.register)
	mux.HandleFunc("/api/v1/auth/login",authHandler.login)
	mux.HandleFunc("/api/v1/auth/me",authHandler.me)
	mux.HandleFunc("/api/v1/auth/logout",authHandler.logout)
	return securityHeaders(mux)
}

func (r *Router) health(w http.ResponseWriter,_ *http.Request){
	writeJSON(w,http.StatusOK,map[string]any{"status":"ok","service":"worldbuster","time":time.Now().UTC()})
}
func (r *Router) worldState(w http.ResponseWriter,_ *http.Request){writeJSON(w,http.StatusOK,r.world.Snapshot())}

func writeJSON(w http.ResponseWriter,status int,value any){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("X-Content-Type-Options","nosniff")
		w.Header().Set("X-Frame-Options","DENY")
		w.Header().Set("Referrer-Policy","strict-origin-when-cross-origin")
		next.ServeHTTP(w,r)
	})
}
