package world

import "testing"

func TestServiceTicks(t *testing.T) {
	s := NewService()
	if s.Snapshot().Tick != 0 {
		t.Fatal("initial tick must be zero")
	}
	s.Tick()
	if s.Snapshot().Tick != 1 {
		t.Fatal("tick did not advance")
	}
}
