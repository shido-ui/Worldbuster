package events

import("crypto/rand";"encoding/hex";"sync";"time")
type Service struct{mu sync.RWMutex;events []Event;notifications map[string][]Notification}
func NewService()*Service{return &Service{notifications:map[string][]Notification{}}}
func id()string{b:=make([]byte,12);_,_=rand.Read(b);return hex.EncodeToString(b)}
func(s *Service)Publish(t,actor,target string,payload map[string]any)Event{s.mu.Lock();defer s.mu.Unlock();e:=Event{ID:id(),Type:t,ActorID:actor,TargetID:target,Payload:payload,CreatedAt:time.Now().UTC()};s.events=append(s.events,e);return e}
func(s *Service)Notify(characterID,t,title,body,eventID string)Notification{s.mu.Lock();defer s.mu.Unlock();n:=Notification{ID:id(),CharacterID:characterID,Type:t,Title:title,Body:body,EventID:eventID,CreatedAt:time.Now().UTC()};s.notifications[characterID]=append(s.notifications[characterID],n);return n}
func(s *Service)Recent(limit int)[]Event{s.mu.RLock();defer s.mu.RUnlock();if limit<=0||limit>100{limit=50};start:=len(s.events)-limit;if start<0{start=0};return append([]Event(nil),s.events[start:]...)}
func(s *Service)Notifications(characterID string)[]Notification{s.mu.RLock();defer s.mu.RUnlock();return append([]Notification(nil),s.notifications[characterID]...)}
