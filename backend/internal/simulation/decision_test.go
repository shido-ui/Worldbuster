package simulation

import "testing"

func TestChooseActionUsesValidatedContext(t *testing.T) {
	c := SimCharacter{Goals: []Goal{{Kind: "WORK", Priority: 5}, {Kind: "SOCIAL", Priority: 3}}}
	a, err := ChooseAction(c, Context{HasJob: true, SocialOpportunity: true, SocialTargetID: "npc2"})
	if err != nil || a.Type != ActionWork {
		t.Fatalf("unexpected action: %+v err=%v", a, err)
	}
	a, err = ChooseAction(c, Context{HasJob: false, SocialOpportunity: true, SocialTargetID: "npc2"})
	if err != nil || a.Type != ActionSocialize {
		t.Fatalf("expected social action: %+v err=%v", a, err)
	}
}

func TestNoActionCases(t *testing.T) {
	c := SimCharacter{Goals: []Goal{{Kind: "WORK", Priority: 10}, {Kind: "STUDY", Priority: 9}, {Kind: "SOCIAL", Priority: 8}, {Kind: "TRAVEL", Priority: 7}}}
	cases := []Context{
		{},
		{HasJob: false},
		{CanStudy: false},
		{SocialOpportunity: true, SocialTargetID: ""},
		{HasTravelDestination: false},
	}
	for i, ctx := range cases {
		if _, err := ChooseAction(c, ctx); err != ErrNoAction {
			t.Fatalf("case %d expected ErrNoAction, got %v", i, err)
		}
	}
}

func TestRestAvailableWhenCritical(t *testing.T) {
	c := SimCharacter{Goals: []Goal{{Kind: "WORK", Priority: 1000}}}
	for _, s := range []BehavioralState{{Needs: Needs{Energy: 10, Rest: 50}}, {Needs: Needs{Energy: 50, Rest: 10}}} {
		a, err := ChooseActionWithState(c, Context{HasJob: true}, s)
		if err != nil || a.Type != ActionRest {
			t.Fatalf("critical rest must dominate: %+v %v", a, err)
		}
	}
	a, err := ChooseActionWithState(c, Context{HasJob: true, RestNeeded: true}, BehavioralState{Needs: Needs{Energy: 100, Rest: 100}})
	if err != nil || a.Type != ActionRest {
		t.Fatalf("RestNeeded must dominate: %+v %v", a, err)
	}
	a, err = ChooseActionWithState(c, Context{HasJob: true}, BehavioralState{Needs: Needs{Energy: 100, Rest: 100}})
	if err != nil || a.Type != ActionWork {
		t.Fatalf("normal state should permit work: %+v %v", a, err)
	}
}

func TestSocialPressureRequiresRealTarget(t *testing.T) {
	c := SimCharacter{Goals: []Goal{{Kind: "SOCIAL", Priority: 50}}, Personality: Personality{Sociability: 100}}
	a, err := ChooseActionWithState(c, Context{SocialOpportunity: true, SocialTargetID: "npc2"}, BehavioralState{Needs: Needs{Energy: 100, Rest: 100, Social: 0, Satisfaction: 50}})
	if err != nil || a.Type != ActionSocialize || a.TargetID != "npc2" {
		t.Fatalf("social pressure failed: %+v %v", a, err)
	}
	if _, err = ChooseAction(c, Context{SocialOpportunity: true}); err != ErrNoAction {
		t.Fatalf("social action without target must be unavailable: %v", err)
	}
}
