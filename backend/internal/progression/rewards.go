package progression

type RewardType string

const (
 RewardAbility RewardType = "ABILITY"
 RewardContent RewardType = "CONTENT"
 RewardStat RewardType = "STAT"
)

type ProgressionReward struct {
 ID string `json:"id"`
 SkillID string `json:"skillId"`
 RequiredLevel int `json:"requiredLevel"`
 Type RewardType `json:"type"`
 Value string `json:"value"`
 Amount int `json:"amount"`
}

func EligibleRewards(state SkillState,rewards []ProgressionReward)[]ProgressionReward{
 out:=make([]ProgressionReward,0)
 for _,r:=range rewards{
  if r.SkillID==state.SkillID && r.RequiredLevel>0 && state.Level>=r.RequiredLevel && r.ID!="" {
   out=append(out,r)
  }
 }
 return out
}
