package progression

import "testing"

func TestProgressionXP(t *testing.T) {
	p := NewProgression()
	if p.Level != 1 {
		t.Fatal("wrong initial level")
	}
	if p.AddXP(100) != 1 || p.Level != 2 {
		t.Fatal("level progression failed")
	}
	if p.XP != 0 {
		t.Fatal("xp carry incorrect")
	}
}
