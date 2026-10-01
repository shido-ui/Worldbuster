package simulation

import (
	"testing"
	"time"

	"github.com/shido-ui/Worldbuster/backend/internal/job"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
	"github.com/shido-ui/Worldbuster/backend/internal/social"
	"github.com/shido-ui/Worldbuster/backend/internal/world"
)

func newLiveAdapterForTest() *LiveAdapter {
	jobs := job.NewService()
	return &LiveAdapter{
		Jobs:        jobs,
		Progression: progression.NewService(),
		Social:      social.NewService(),
		Travel:      world.NewTravelService(),
		Locations:   map[string]string{"c1": "central", "c2": "harbor"},
	}
}

func TestLiveAdapterRestExecution(t *testing.T) {
	a := newLiveAdapterForTest()
	result := a.Execute(Action{Type: ActionRest}, SimCharacter{ID: "c1"})
	if !result.Accepted {
		t.Fatalf("rest rejected: %s", result.Reason)
	}
}

func TestLiveAdapterTravelExecution(t *testing.T) {
	a := newLiveAdapterForTest()
	result := a.Execute(Action{Type: ActionTravel, TargetID: "harbor"}, SimCharacter{ID: "c1"})
	if !result.Accepted {
		t.Fatalf("travel rejected: %s", result.Reason)
	}
	state, ok := a.Travel.Get("c1", time.Now())
	if !ok {
		t.Fatal("accepted travel was not recorded")
	}
	if state.From != "central" || state.To != "harbor" {
		t.Fatalf("unexpected travel state: %+v", state)
	}
}

func TestLiveAdapterRejectsUnknownTravelDestination(t *testing.T) {
	a := newLiveAdapterForTest()
	result := a.Execute(Action{Type: ActionTravel, TargetID: "unknown"}, SimCharacter{ID: "c1"})
	if result.Accepted {
		t.Fatal("unknown travel destination was accepted")
	}
}

func TestLiveAdapterContextExposesTravelOpportunity(t *testing.T) {
	a := newLiveAdapterForTest()
	ctx := a.Context("c1")
	if !ctx.HasTravelDestination {
		t.Fatal("context failed to expose a reachable travel destination")
	}
	a.Locations["c1"] = "highlands"
	ctx = a.Context("c1")
	if !ctx.HasTravelDestination {
		t.Fatal("highlands should expose its oldtown route")
	}
}
