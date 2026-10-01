package simulation

import "sync"

type GroupRole string

const (
	GroupLeader  GroupRole = "LEADER"
	GroupOfficer GroupRole = "OFFICER"
	GroupMember  GroupRole = "MEMBER"
)

type GroupResource struct {
	Currency  int64 `json:"currency"`
	Influence int   `json:"influence"`
}

type GroupRuntime struct {
	mu        sync.RWMutex
	roles     map[string]map[string]GroupRole
	resources map[string]GroupResource
	goals     map[string]string
}

func NewGroupRuntime() *GroupRuntime {
	return &GroupRuntime{roles: map[string]map[string]GroupRole{}, resources: map[string]GroupResource{}, goals: map[string]string{}}
}

func (r *GroupRuntime) SetRole(groupID, characterID string, role GroupRole) bool {
	if groupID == "" || characterID == "" {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.roles[groupID] == nil {
		r.roles[groupID] = map[string]GroupRole{}
	}
	r.roles[groupID][characterID] = role
	return true
}

func (r *GroupRuntime) Role(groupID, characterID string) (GroupRole, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := r.roles[groupID]
	if m == nil {
		return "", false
	}
	v, ok := m[characterID]
	return v, ok
}

func (r *GroupRuntime) AddResource(groupID string, currency int64, influence int) bool {
	if groupID == "" {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	v := r.resources[groupID]
	v.Currency += currency
	v.Influence += influence
	if v.Currency < 0 {
		v.Currency = 0
	}
	if v.Influence < 0 {
		v.Influence = 0
	}
	r.resources[groupID] = v
	return true
}

func (r *GroupRuntime) Resource(groupID string) (GroupResource, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.resources[groupID]
	return v, ok
}

func (r *GroupRuntime) SetGoal(groupID, goal string) bool {
	if groupID == "" || goal == "" {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.goals[groupID] = goal
	return true
}

func (r *GroupRuntime) Goal(groupID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.goals[groupID]
	return v, ok
}

func (r *GroupRuntime) ApplyDecision(groupID string, action GroupAction) bool {
	if groupID == "" || action.Type == "" {
		return false
	}
	ApplyGroupAction(r, groupID, action)
	return true
}
