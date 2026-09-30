package player

import (
	"errors"
	"sync"
)

var ErrInvalidProfile = errors.New("invalid player profile")

type Service struct {
	mu      sync.RWMutex
	profile map[string]Profile
}

func NewService() *Service {
	return &Service{profile: make(map[string]Profile)}
}

func (s *Service) Get(id string) (Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profile[id]
	return p, ok
}

func (s *Service) Create(p Profile) error {
	if p.ID == "" || p.AccountID == "" || p.DisplayName == "" {
		return ErrInvalidProfile
	}
	if p.Level < 1 || p.XP < 0 || p.Cash < 0 || p.Energy < 0 || p.Energy > 100 {
		return ErrInvalidProfile
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.profile[p.ID]; exists {
		return errors.New("player already exists")
	}
	s.profile[p.ID] = p
	return nil
}
