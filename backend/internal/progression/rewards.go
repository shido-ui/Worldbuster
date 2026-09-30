package progression

type RewardType string

const (
 RewardAbility RewardType = "ABILITY"
 RewardContent RewardType = "CONTENT"
 RewardStat RewardType = "STAT"
)

type ProgressionReward struct {
 ID string
 SkillID string
 RequiredLevel int
 Type RewardType
 Value string
 Amount int
}

func EligibleRewards(state SkillState, rewards []ProgressionReward) []ProgressionReward {
 out:=make([]ProgressionReward,0)
 for _,r:=range rewards {
  if r.ID!="" && r.SkillID==state.SkillID && r.RequiredLevel>0 && state.Level>=r.RequiredLevel {
   out=append(out,r)
  }
 }
 return out
}
