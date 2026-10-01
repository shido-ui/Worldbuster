package organization

import (
	"errors"
	"sync"
)

var ErrInvalidOrganization = errors.New("invalid organization")
var ErrAlreadyMember = errors.New("already a member")
var ErrNotMember = errors.New("not a member")
var ErrFull = errors.New("organization is full")

type Service struct {
	mu      sync.RWMutex
	orgs    map[string]*Organization
	members map[string]map[string]Membership
}

func NewService() *Service {
	return &Service{orgs: map[string]*Organization{}, members: map[string]map[string]Membership{}}
}
func (s *Service) Create(o Organization) error {
	if o.ID == "" || o.Name == "" || o.MaxMembers <= 0 {
		return ErrInvalidOrganization
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.orgs[o.ID]; ok {
		return ErrInvalidOrganization
	}
	s.orgs[o.ID] = &o
	s.members[o.ID] = map[string]Membership{}
	return nil
}
func (s *Service) Join(orgID, charID, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orgs[orgID]
	if !ok {
		return ErrInvalidOrganization
	}
	m := s.members[orgID]
	if _, ok := m[charID]; ok {
		return ErrAlreadyMember
	}
	if len(m) >= o.MaxMembers {
		return ErrFull
	}
	m[charID] = Membership{OrganizationID: orgID, CharacterID: charID, Role: role}
	return nil
}
func (s *Service) Leave(orgID, charID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.members[orgID]
	if !ok {
		return ErrInvalidOrganization
	}
	if _, ok := m[charID]; !ok {
		return ErrNotMember
	}
	delete(m, charID)
	return nil
}
func (s *Service) Get(orgID string) (Organization, []Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orgs[orgID]
	if !ok {
		return Organization{}, nil, ErrInvalidOrganization
	}
	out := make([]Membership, 0, len(s.members[orgID]))
	for _, m := range s.members[orgID] {
		out = append(out, m)
	}
	return *o, out, nil
}
