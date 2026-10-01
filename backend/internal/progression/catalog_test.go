package progression

import "testing"

func TestDefaultSkillCatalog(t *testing.T) {
	c := NewDefaultSkillCatalog()
	if c == nil || len(c.List()) < 4 {
		t.Fatal("default catalog incomplete")
	}
}
