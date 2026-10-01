package simulation

import (
	"testing"
	"time"
)

func TestUpdateBehavior(t *testing.T) {
	s := BehavioralState{Needs: Needs{Energy: 50, Satisfaction: 50}}
	UpdateBehavior(&s, ActionWork, time.Unix(1, 0))
	if s.Needs.Energy >= 50 || s.Stress <= 0 || len(s.Memories) != 1 {
		t.Fatal("behavior did not evolve")
	}
	UpdateBehavior(&s, ActionRest, time.Unix(2, 0))
	if s.Needs.Energy <= 0 {
		t.Fatal("rest did not restore energy")
	}
}
