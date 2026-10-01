package events

import "testing"

func TestNewsGenerator(t *testing.T) {
	n := NewsGenerator{}.Generate(WorldEvent{ID: "e1", Summary: "A major group event", Importance: 90})
	if n.EventID != "e1" || n.Headline != "A major group event" {
		t.Fatal("news generation failed")
	}
}
