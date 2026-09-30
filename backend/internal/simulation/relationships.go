package simulation

import "sync"

type RelationshipKind string
const (
 RelationshipFriend RelationshipKind = "FRIEND"
 RelationshipRival RelationshipKind = "RIVAL"
 RelationshipAcquaintance RelationshipKind = "ACQUAINTANCE"
)

type Relationship struct {
 FromID string `json:"fromId"`
 ToID string `json:"toId"`
 Kind RelationshipKind `json:"kind"`
 Familiarity int `json:"familiarity"`
 Trust int `json:"trust"`
}

type RelationshipService struct {
 mu sync.RWMutex
 relationships map[string]Relationship
}

func NewRelationshipService()*RelationshipService{return &RelationshipService{relationships:map[string]Relationship{}}}
func relationshipKey(a,b string)string{return a+":"+b}

func(s *RelationshipService) Get(a,b string)(Relationship,bool){s.mu.RLock();defer s.mu.RUnlock();r,ok:=s.relationships[relationshipKey(a,b)];return r,ok}

func(s *RelationshipService) Update(a,b string, kind RelationshipKind, familiarityDelta,trustDelta int){
 if a==""||b==""||a==b{return}
 s.mu.Lock();defer s.mu.Unlock()
 r:=s.relationships[relationshipKey(a,b)]
 r.FromID=a;r.ToID=b;r.Kind=kind
 r.Familiarity=clamp(r.Familiarity+familiarityDelta)
 r.Trust=clampSigned(r.Trust+trustDelta)
 s.relationships[relationshipKey(a,b)]=r
}

func clamp(v int)int{if v<0{return 0};if v>100{return 100};return v}
func clampSigned(v int)int{if v< -100{return -100};if v>100{return 100};return v}
