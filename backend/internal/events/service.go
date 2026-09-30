package events

import("crypto/rand";"encoding/hex";"errors";"sync";"time")
type Service struct{mu sync.RWMutex;events []Event;notifications map[string][]Notification;subscribers map[chan Event]struct{}}
func NewService()*Service{return &Service{notifications:map[string][]Notification{},subscribers:map[chan Event]struct{}{}}}
func id()(string,error){b:=make([]byte,12);if _,err:=rand.Read(b);err!=nil{return "",err};return hex.EncodeToString(b),nil}
func(s *Service)Publish(t,actor,target string,payload map[string]any)(Event,error){s.mu.Lock();defer s.mu.Unlock();eventID,err:=id();if err!=nil{return Event{},errors.New("generate event id: "+err.Error())};e:=Event{ID:eventID,Type:t,ActorID:actor,TargetID:target,Payload:payload,CreatedAt:time.Now().UTC()};for ch:=range s.subscribers{select{case ch<-e:default:}};s.events=append(s.events,e);if len(s.events)>1000{s.events=s.events[len(s.events)-1000:]};return e}
func(s *Service)Notify(characterID,t,title,body,eventID string)(Notification,error){s.mu.Lock();defer s.mu.Unlock();notificationID,err:=id();if err!=nil{return Notification{},errors.New("generate notification id: "+err.Error())};n:=Notification{ID:notificationID,CharacterID:characterID,Type:t,Title:title,Body:body,EventID:eventID,CreatedAt:time.Now().UTC()};s.notifications[characterID]=append(s.notifications[characterID],n);return n}
func(s *Service)Recent(limit int)[]Event{s.mu.RLock();defer s.mu.RUnlock();if limit<=0||limit>100{limit=50};start:=len(s.events)-limit;if start<0{start=0};return append([]Event(nil),s.events[start:]...)}
func(s *Service)Notifications(characterID string)[]Notification{s.mu.RLock();defer s.mu.RUnlock();return append([]Notification(nil),s.notifications[characterID]...)}
func(s *Service)MarkRead(characterID,notificationID string)bool{s.mu.Lock();defer s.mu.Unlock();for i:=range s.notifications[characterID]{if s.notifications[characterID][i].ID==notificationID{s.notifications[characterID][i].Read=true;return true}};return false}
func(s *Service)MarkAllRead(characterID string)int{s.mu.Lock();defer s.mu.Unlock();n:=0;for i:=range s.notifications[characterID]{if !s.notifications[characterID][i].Read{s.notifications[characterID][i].Read=true;n++}};return n}

func(s *Service)Subscribe()chan Event{s.mu.Lock();defer s.mu.Unlock();ch:=make(chan Event,16);s.subscribers[ch]=struct{}{};return ch}
func(s *Service)Unsubscribe(ch chan Event){s.mu.Lock();defer s.mu.Unlock();if _,ok:=s.subscribers[ch];ok{delete(s.subscribers,ch);close(ch)}}