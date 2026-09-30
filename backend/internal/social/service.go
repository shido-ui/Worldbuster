package social

import("errors";"sync";"time";"crypto/rand";"encoding/hex")
var ErrInvalidMessage=errors.New("invalid message")
var ErrSelfTarget=errors.New("self target is not allowed")
type Service struct{mu sync.RWMutex;relationships map[string]Relationship;messages map[string][]Message}
func NewService()*Service{return &Service{relationships:map[string]Relationship{},messages:map[string][]Message{}}}
func key(a,b string)string{return a+":"+b}
func(s *Service)SetRelationship(a,b,kind string,score int)error{if a==""||b==""||a==b{return ErrSelfTarget};s.mu.Lock();defer s.mu.Unlock();s.relationships[key(a,b)]=Relationship{CharacterID:a,TargetID:b,Kind:kind,Score:score};return nil}
func(s *Service)Relationship(a,b string)(Relationship,bool){s.mu.RLock();defer s.mu.RUnlock();r,ok:=s.relationships[key(a,b)];return r,ok}
func(s *Service)Send(from,to,body string)(Message,error){if from==""||to==""||from==to||body==""{return Message{},ErrInvalidMessage};b:=make([]byte,12);_,_=rand.Read(b);m:=Message{ID:hex.EncodeToString(b),FromID:from,ToID:to,Body:body,CreatedAt:time.Now().Unix()};s.mu.Lock();defer s.mu.Unlock();s.messages[to]=append(s.messages[to],m);return m,nil}
func(s *Service)Inbox(id string)[]Message{s.mu.RLock();defer s.mu.RUnlock();return append([]Message(nil),s.messages[id]...)}
