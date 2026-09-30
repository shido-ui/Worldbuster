package simulation

import (
 "sync"
 "time"
)

type StateService struct {
 mu sync.RWMutex
 states map[string]*BehavioralState
}

func NewStateService()*StateService{return &StateService{states:map[string]*BehavioralState{}}}

func(s *StateService) Get(id string) BehavioralState {
 s.mu.RLock(); defer s.mu.RUnlock()
 if v,ok:=s.states[id];ok{return cloneState(*v)}
 return BehavioralState{Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50},Mood:25}
}

func(s *StateService) Tick(id string){
 s.mu.Lock(); defer s.mu.Unlock()
 st:=s.getLocked(id)
 // Needs drift slowly between actions, creating pressure for future decisions.
 st.Needs.Energy-=1
 st.Needs.Social-=1
 st.Needs.Rest-=1
 if st.Needs.Satisfaction>0 {st.Needs.Satisfaction-=1}
 clampNeeds(&st.Needs)
 if st.Stress>100 {st.Stress=100}
 if st.Stress>0 {st.Stress--}
 st.Mood=st.Needs.Satisfaction-st.Stress/2
}

func(s *StateService) Apply(id string, action ActionType, now time.Time){
 s.mu.Lock(); defer s.mu.Unlock()
 st:=s.getLocked(id)
 UpdateBehavior(st,action,now)
}

func(s *StateService) getLocked(id string)*BehavioralState{
 if st,ok:=s.states[id];ok{return st}
 st:=&BehavioralState{Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50},Mood:25}
 s.states[id]=st
 return st
}

func cloneState(st BehavioralState)BehavioralState{
 st.Memories=append([]Memory(nil),st.Memories...)
 return st
}
