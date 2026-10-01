package simulation

import (
	"testing"
	"time"
)

func TestSocialRuntime(t *testing.T) {
	s := &SocialRuntime{Relationships: NewRelationshipService()}
	s.ApplySocialAction("a", "b", ActionSocialize, time.Unix(1, 0))
	r, ok := s.Relationships.Get("a", "b")
	if !ok || r.Familiarity != 8 {
		t.Fatal("social action did not create relationship")
	}
	s.ApplySocialAction("a", "b", ActionSocialize, time.Unix(2, 0))
	r, _ = s.Relationships.Get("a", "b")
	if r.Kind != RelationshipFriend {
		t.Fatal("repeated social interaction did not deepen relationship")
	}
}
func TestSocialDecisionBias(t *testing.T) {
	if SocialDecisionBias(Relationship{Kind: RelationshipFriend, Familiarity: 100, Trust: 100}) <= SocialDecisionBias(Relationship{Kind: RelationshipRival, Familiarity: 100, Trust: 100}) {
		t.Fatal("friend should influence social choice more positively")
	}
}
