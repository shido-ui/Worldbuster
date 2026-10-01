package simulation

import "testing"

func TestDeterministicBusinessID(t *testing.T) {
	first := DeterministicBusinessID("npc-1", BusinessProduction)
	second := DeterministicBusinessID("npc-1", BusinessProduction)
	if first != second {
		t.Fatalf("business ID is not deterministic: %q != %q", first, second)
	}
	if first == DeterministicBusinessID("npc-2", BusinessProduction) {
		t.Fatal("different owners produced the same business ID")
	}
	if first == DeterministicBusinessID("npc-1", BusinessRetail) {
		t.Fatal("different business types produced the same business ID")
	}
}