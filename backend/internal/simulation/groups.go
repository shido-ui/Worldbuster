package simulation

import "sync"

type Group struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	LeaderID  string   `json:"leaderId"`
	MemberIDs []string `json:"memberIds"`
}

type GroupService struct {
	mu         sync.RWMutex
	groups     map[string]*Group
	membership map[string]string
}

func NewGroupService() *GroupService {
	return &GroupService{groups: map[string]*Group{}, membership: map[string]string{}}
}

func (s *GroupService) Create(id, name, leader string) bool {
	if id == "" || name == "" || leader == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[id]; ok {
		return false
	}
	if _, ok := s.membership[leader]; ok {
		return false
	}
	s.groups[id] = &Group{ID: id, Name: name, LeaderID: leader, MemberIDs: []string{leader}}
	s.membership[leader] = id
	return true
}

func (s *GroupService) Join(groupID, characterID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.groups[groupID]
	if !ok || characterID == "" {
		return false
	}
	if _, ok := s.membership[characterID]; ok {
		return false
	}
	g.MemberIDs = append(g.MemberIDs, characterID)
	s.membership[characterID] = groupID
	return true
}

func (s *GroupService) Get(id string) (Group, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[id]
	if !ok {
		return Group{}, false
	}
	out := *g
	out.MemberIDs = append([]string(nil), g.MemberIDs...)
	return out, true
}

func (s *GroupService) GroupOf(characterID string) (Group, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.membership[characterID]
	if !ok {
		return Group{}, false
	}
	g := s.groups[id]
	out := *g
	out.MemberIDs = append([]string(nil), g.MemberIDs...)
	return out, true
}

func (s *GroupService) List() []Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Group, 0, len(s.groups))
	for _, g := range s.groups {
		cp := *g
		cp.MemberIDs = append([]string(nil), g.MemberIDs...)
		out = append(out, cp)
	}
	return out
}
