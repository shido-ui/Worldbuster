package simulation

import ("sync";"time")

type StateService struct{mu sync.RWMutex;states map[string]BehavioralState}
func NewStateService()*StateService{return &StateService{states:map[string]BehavioralState{}}}
func(s *StateService)Get(id string)BehavioralState{s.mu.RLock();defer s.mu.RUnlock();x,ok:=s.states[id];if !ok{return defaultState()};return cloneState(x)}
func defaultState()BehavioralState{return BehavioralState{Mood:25,Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50},Memories:[]Memory{},Relationships:map[string]Relationship{}}}
func cloneState(x BehavioralState)BehavioralState{y:=x;y.Memories=append([]Memory(nil),x.Memories...);y.Relationships=map[string]Relationship{};for k,v:=range x.Relationships{y.Relationships[k]=v};return y}
func(s *StateService)ensure(id string)BehavioralState{if x,ok:=s.states[id];ok{return x};x:=defaultState();s.states[id]=x;return x}
func(s *StateService)Apply(id string,a ActionType,now time.Time){s.mu.Lock();defer s.mu.Unlock();x:=s.ensure(id);x=applyBehavior(x,a,now);s.states[id]=x}
func(s *StateService)Tick(id string){s.mu.Lock();defer s.mu.Unlock();x:=s.ensure(id);x.Needs.Energy=clamp(x.Needs.Energy-1,0,100);x.Needs.Rest=clamp(x.Needs.Rest-1,0,100);x.Needs.Social=clamp(x.Needs.Social-1,0,100);s.states[id]=x}
func(s *StateService)RecordRelationship(id string,r Relationship,now time.Time){s.mu.Lock();defer s.mu.Unlock();x:=s.ensure(id);if x.Relationships==nil{x.Relationships=map[string]Relationship{}};x.Relationships[r.ToID]=r;x.Memories=append(x.Memories,Memory{Type:"SOCIAL_RELATIONSHIP",TargetID:r.ToID,Event:"social interaction",Importance:3,Sentiment:r.Affinity,CreatedAt:now});if len(x.Memories)>50{x.Memories=x.Memories[len(x.Memories)-50:]};s.states[id]=x}
