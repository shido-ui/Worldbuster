package progression

import "testing"

func TestApplyXPWithRewards(t *testing.T) {
	ps := NewPlayerProgressionStore()
	us := NewUnlockStore()
	rewards := []ProgressionReward{{ID: "r1", SkillID: "focus", RequiredLevel: 5, Type: RewardAbility, Value: "advanced_focus"}}
	result, err := ApplyXPWithRewards(ps, us, "p1", 100, SkillState{SkillID: "focus", Level: 5}, rewards)
	if err != nil || result.Progression.Level != 2 {
		t.Fatal("progression transaction failed")
	}
	if len(result.Granted) != 0 {
		t.Fatal("reward should not grant before milestone")
	}
}
