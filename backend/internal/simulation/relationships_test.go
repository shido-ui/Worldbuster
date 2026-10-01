package simulation

import "testing"

func TestRelationshipService(t *testing.T) {
	s := NewRelationshipService()
	s.Update("a", "b", RelationshipAcquaintance, 20, 10)
	r, ok := s.Get("a", "b")
	if !ok || r.Familiarity != 20 || r.Trust != 10 {
		t.Fatal("relationship not stored")
	}
	s.Update("a", "b", RelationshipFriend, 100, 100)
	r, _ = s.Get("a", "b")
	if r.Familiarity != 100 || r.Trust != 100 {
		t.Fatal("relationship not clamped")
	}
}
