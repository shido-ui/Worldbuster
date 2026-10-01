package simulation

import (
	"sync"
	"time"
)

type StateService struct {
	mu     sync.RWMutex
	states map[string]BehavioralState
}

func NewStateService() *StateService { return &StateService{states: map[string]BehavioralState{}} }
func (s *StateService) Get(id string) BehavioralState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.states[id]
	if !ok {
		return defaultState()
	}
	return cloneState(x)
}
func defaultState() BehavioralState {
	return BehavioralState{Mood: 25, Needs: Needs{Energy: 100, Social: 50, Rest: 50, Satisfaction: 50}, Memories: []Memory{}, Relationships: map[string]Relationship{}, Economy: EconomicState{Balance: 100}}
}
func cloneState(x BehavioralState) BehavioralState {
	y := x
	y.Memories = append([]Memory(nil), x.Memories...)
	y.Relationships = map[string]Relationship{}
	for k, v := range x.Relationships {
		y.Relationships[k] = v
	}
	return y
}
func (s *StateService) ensure(id string) BehavioralState {
	if x, ok := s.states[id]; ok {
		return x
	}
	x := defaultState()
	s.states[id] = x
	return x
}
func (s *StateService) Apply(id string, a ActionType, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	x := s.ensure(id)
	UpdateBehavior(&x, a, now)
	s.states[id] = x
}
func (s *StateService) Tick(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	x := s.ensure(id)
	x.Needs.Energy = clamp(x.Needs.Energy - 1)
	x.Needs.Rest = clamp(x.Needs.Rest - 1)
	x.Needs.Social = clamp(x.Needs.Social - 1)
	s.states[id] = x
}
func (s *StateService) RecordRelationship(id string, r Relationship, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	x := s.ensure(id)
	if x.Relationships == nil {
		x.Relationships = map[string]Relationship{}
	}
	x.Relationships[r.ToID] = r
	x.Memories = append(x.Memories, Memory{Type: "SOCIAL_RELATIONSHIP", TargetID: r.ToID, Event: "social interaction", Importance: 3, Sentiment: r.Affinity, CreatedAt: now.Unix()})
	if len(x.Memories) > 50 {
		x.Memories = x.Memories[len(x.Memories)-50:]
	}
	s.states[id] = x
}
func (s *StateService) CanAffordAction(id string, a ActionType) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.states[id]
	if !ok {
		x = defaultState()
	}
	cost := economicCost(a)
	return cost == 0 || x.Economy.CanAfford(cost)
}
func (s *StateService) ApplyEconomy(id string, a ActionType) (EconomicState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	x := s.ensure(id)
	if cost := economicCost(a); cost > 0 && !x.Economy.Spend(cost) {
		return x.Economy, false
	}
	if a == ActionWork {
		x.Economy.Earn(10)
	}
	s.states[id] = x
	return x.Economy, true
}

func (s *StateService) Remember(id, kind, target, event string, importance, sentiment int) Memory {
	s.mu.Lock()
	defer s.mu.Unlock()
	x := s.ensure(id)
	m := Memory{ID: time.Now().UTC().Format("20060102150405.000000000"), Type: kind, TargetID: target, Event: event, Importance: importance, Sentiment: sentiment, CreatedAt: time.Now().UTC().Unix()}
	x.Memories = append(x.Memories, m)
	if len(x.Memories) > 50 {
		x.Memories = x.Memories[len(x.Memories)-50:]
	}
	s.states[id] = x
	return m
}
