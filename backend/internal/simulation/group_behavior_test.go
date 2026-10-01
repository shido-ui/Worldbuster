package simulation

import "testing"

func TestChooseGroupAction(t *testing.T) {
	a := ChooseGroupAction(GroupContext{MemberCount: 5, Currency: 100, Goal: "grow", NearbyCandidates: 3})
	if a.Type != GroupInvest {
		t.Fatalf("expected investment, got %s", a.Type)
	}
}

func TestApplyGroupAction(t *testing.T) {
	r := NewGroupRuntime()
	r.AddResource("g", 100, 0)
	ApplyGroupAction(r, "g", GroupInvest)
	v, _ := r.Resource("g")
	if v.Currency != 50 || v.Influence != 10 {
		t.Fatal("investment consequence incorrect")
	}
}
