package progression

import "testing"

func TestEligibleRewards(t *testing.T) {
	rewards := []ProgressionReward{{ID: "r1", SkillID: "focus", RequiredLevel: 5, Type: RewardAbility, Value: "advanced_focus"}}
	got := EligibleRewards(SkillState{SkillID: "focus", Level: 5}, rewards)
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatal("reward not eligible")
	}
}
