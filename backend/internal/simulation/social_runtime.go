package simulation

import "time"

type SocialRuntime struct {
 Relationships *RelationshipService
}

func (s *SocialRuntime) ApplySocialAction(actor, target string, action ActionType, now time.Time) {
 if s==nil||s.Relationships==nil||actor==""||target==""||actor==target{return}
 switch action {
 case ActionSocialize:
  r,ok:=s.Relationships.Get(actor,target)
  if !ok {
   s.Relationships.Update(actor,target,RelationshipAcquaintance,8,2)
   return
  }
  if r.Kind==RelationshipRival {
   s.Relationships.Update(actor,target,RelationshipRival,2,-1)
  } else {
   s.Relationships.Update(actor,target,RelationshipFriend,6,3)
  }
 }
}

func SocialDecisionBias(rel Relationship) int {
 familiarity:=rel.Familiarity/5
 trust:=rel.Trust/10
 switch rel.Kind {
 case RelationshipFriend:
  return familiarity+trust+10
 case RelationshipRival:
  return -familiarity-trust
 default:
  return familiarity+trust
 }
}
