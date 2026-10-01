package progression

import "fmt"

type TransactionResult struct {
	Progression Progression
	Granted     []ProgressionReward
}

func ApplyXPWithRewards(store *PlayerProgressionStore, unlocks *UnlockStore, playerID string, amount int64, skillState SkillState, rewards []ProgressionReward) (TransactionResult, error) {
	if store == nil || unlocks == nil || playerID == "" || amount <= 0 {
		return TransactionResult{}, fmt.Errorf("invalid progression transaction")
	}
	p, gained := store.AddXP(playerID, amount)
	if gained == 0 {
		return TransactionResult{Progression: p}, nil
	}
	eligible := EligibleRewards(skillState, rewards)
	granted := make([]ProgressionReward, 0, len(eligible))
	for _, r := range eligible {
		if p.Level >= r.RequiredLevel && unlocks.Grant(playerID, r) {
			granted = append(granted, r)
		}
	}
	return TransactionResult{Progression: p, Granted: granted}, nil
}
