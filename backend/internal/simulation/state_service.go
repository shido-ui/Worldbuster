package simulation
import("crypto/rand";"encoding/hex";"sync";"time")
type StateService struct{mu sync.RWMutex;states map[string]*BehavioralState}
func NewStateService()*StateService{return &StateService{states:map[string]*BehavioralState{}}}
func(s *StateService)Get(id string)BehavioralState{s.mu.Lock();defer s.mu.Unlock();st,ok:=s.states[id];if !ok{st=&BehavioralState{Mood:50,Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50}};s.states[id]=st};return cloneState(*st)}
func(s *StateService)Tick(id string)BehavioralState{s.mu.Lock();defer s.mu.Unlock();st,ok:=s.states[id];if !ok{st=&BehavioralState{Mood:50,Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50}};s.states[id]=st};if st.Needs.Energy>0{st.Needs.Energy--};if st.Needs.Rest>0{st.Needs.Rest--};if st.Needs.Social>0{st.Needs.Social--};if st.Needs.Energy<25||st.Needs.Rest<20{st.Stress++};return cloneState(*st)}
func(s *StateService)Remember(id,typ,target,event string,importance,sentiment int)Memory{s.mu.Lock();defer s.mu.Unlock();st,ok:=s.states[id];if !ok{st=&BehavioralState{Mood:50,Needs:Needs{Energy:100,Social:50,Rest:50,Satisfaction:50}};s.states[id]=st};b:=make([]byte,8);_,_=rand.Read(b);m:=Memory{ID:hex.EncodeToString(b),Type:typ,TargetID:target,Event:event,Importance:importance,Sentiment:sentiment,CreatedAt:time.Now().Unix()};st.Memories=append(st.Memories,m);if len(st.Memories)>100{st.Memories=st.Memories[len(st.Memories)-100:]};return m}
func cloneState(s BehavioralState)BehavioralState{s.Memories=append([]Memory(nil),s.Memories...);return s}
