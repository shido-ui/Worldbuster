package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrInvalidUsername    = errors.New("username must be 3-24 characters")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrSessionNotFound    = errors.New("session not found")
)

type Service struct {
	mu       sync.RWMutex
	accounts map[string]Account
	byName   map[string]string
	sessions map[string]Session
}

func NewService() *Service {
	return &Service{accounts: make(map[string]Account), byName: make(map[string]string), sessions: make(map[string]Session)}
}

func (s *Service) Register(username, password string) (PublicAccount, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 24 {
		return PublicAccount{}, ErrInvalidUsername
	}
	if len(password) < 8 {
		return PublicAccount{}, ErrInvalidPassword
	}
	key := strings.ToLower(username)

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byName[key]; exists {
		return PublicAccount{}, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return PublicAccount{}, err
	}
	id, err := randomID()
	if err != nil {
		return PublicAccount{}, err
	}
	now := time.Now().UTC()
	account := Account{ID: id, Username: username, PasswordHash: string(hash), CreatedAt: now, UpdatedAt: now}
	s.accounts[id] = account
	s.byName[key] = id
	return account.Public(), nil
}

func (s *Service) Authenticate(username, password string) (PublicAccount, error) {
	s.mu.RLock()
	id, ok := s.byName[strings.ToLower(strings.TrimSpace(username))]
	account := s.accounts[id]
	s.mu.RUnlock()
	if !ok || bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)) != nil {
		return PublicAccount{}, ErrInvalidCredentials
	}
	return account.Public(), nil
}

func (s *Service) CreateSession(accountID string, ttl time.Duration) (Session, error) {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[accountID]; !ok {
		return Session{}, ErrInvalidCredentials
	}
	id, err := randomID()
	if err != nil {
		return Session{}, err
	}
	now := time.Now().UTC()
	session := Session{ID: id, AccountID: accountID, CreatedAt: now, ExpiresAt: now.Add(ttl)}
	s.sessions[id] = session
	return session, nil
}

func (s *Service) ResolveSession(id string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if !ok {
		return Session{}, ErrSessionNotFound
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		delete(s.sessions, id)
		return Session{}, ErrSessionNotFound
	}
	return session, nil
}

func (s *Service) RevokeSession(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

func randomID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
