package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/shido-ui/Worldbuster/backend/internal/auth"
)

type authHandler struct{ service *auth.Service }

type credentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newAuthHandler(service *auth.Service) *authHandler { return &authHandler{service:service} }

func (h *authHandler) register(w http.ResponseWriter,r *http.Request) {
	var req credentialsRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil { writeJSON(w,400,map[string]string{"error":"invalid json"}); return }
	account,err:=h.service.Register(req.Username,req.Password)
	if err!=nil {
		status:=400
		if errors.Is(err,auth.ErrUsernameTaken){status=409}
		writeJSON(w,status,map[string]string{"error":err.Error()}); return
	}
	writeJSON(w,201,account)
}

func (h *authHandler) login(w http.ResponseWriter,r *http.Request) {
	var req credentialsRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil { writeJSON(w,400,map[string]string{"error":"invalid json"}); return }
	account,err:=h.service.Authenticate(req.Username,req.Password)
	if err!=nil { writeJSON(w,401,map[string]string{"error":"invalid credentials"}); return }
	session,err:=h.service.CreateSession(account.ID,24*time.Hour)
	if err!=nil { writeJSON(w,500,map[string]string{"error":"session creation failed"}); return }
	http.SetCookie(w,&http.Cookie{Name:"worldbuster_session",Value:session.ID,Path:"/",HttpOnly:true,SameSite:http.SameSiteLaxMode,Secure:isSecureRequest(r),MaxAge:int(time.Until(session.ExpiresAt).Seconds())})
	writeJSON(w,200,account)
}

func (h *authHandler) me(w http.ResponseWriter,r *http.Request) {
	session,ok:=sessionFromRequest(h.service,r)
	if !ok { writeJSON(w,401,map[string]string{"error":"authentication required"}); return }
	writeJSON(w,200,map[string]string{"accountId":session.AccountID})
}

func (h *authHandler) logout(w http.ResponseWriter,r *http.Request) {
	if c,err:=r.Cookie("worldbuster_session"); err==nil { h.service.RevokeSession(c.Value) }
	http.SetCookie(w,&http.Cookie{Name:"worldbuster_session",Value:"",Path:"/",HttpOnly:true,MaxAge:-1,SameSite:http.SameSiteLaxMode,Secure:isSecureRequest(r)})
	w.WriteHeader(http.StatusNoContent)
}

func sessionFromRequest(service *auth.Service,r *http.Request)(auth.Session,bool) {
	c,err:=r.Cookie("worldbuster_session")
	if err!=nil || strings.TrimSpace(c.Value)=="" { return auth.Session{},false }
	session,err:=service.ResolveSession(c.Value)
	return session,err==nil
}

func isSecureRequest(r *http.Request) bool { return r.TLS!=nil || r.Header.Get("X-Forwarded-Proto")=="https" }
