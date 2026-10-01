package auth

import (
	"testing"
	"time"
)

func TestRegisterAuthenticateAndSession(t *testing.T) {
	s := NewService()
	account, err := s.Register("Shido", "secure-pass-123")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Authenticate("Shido", "wrong-password"); err != ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
	if _, err := s.Authenticate("Shido", "secure-pass-123"); err != nil {
		t.Fatal(err)
	}
	session, err := s.CreateSession(account.ID, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.ResolveSession(session.ID); err != nil {
		t.Fatal(err)
	}
	s.RevokeSession(session.ID)
	if _, err := s.ResolveSession(session.ID); err != ErrSessionNotFound {
		t.Fatalf("expected revoked session, got %v", err)
	}
}
