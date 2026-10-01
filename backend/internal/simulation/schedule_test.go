package simulation

import "testing"

func TestSchedule(t *testing.T) {
	s := Schedule{Blocks: []ScheduleBlock{{StartHour: 9, EndHour: 17, Action: ActionWork}}}
	if a, ok := s.ActionAt(12); !ok || a != ActionWork {
		t.Fatal("schedule failed")
	}
	if _, ok := s.ActionAt(20); ok {
		t.Fatal("unexpected scheduled action")
	}
}
