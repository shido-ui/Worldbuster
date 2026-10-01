package events

import (
	"testing"
	"time"
)

func TestHistory(t *testing.T) {
	h := NewHistory()
	h.Record(WorldEvent{ID: "1", Type: "GROUP_ACTION", Summary: "A group changed", At: time.Now(), Importance: 80})
	h.Record(WorldEvent{ID: "2", Type: "MINOR", Summary: "Minor event", At: time.Now(), Importance: 20})
	if len(h.Recent(10)) != 2 {
		t.Fatal("recent history missing")
	}
	if len(h.Important(10)) != 1 {
		t.Fatal("important history missing")
	}
}
