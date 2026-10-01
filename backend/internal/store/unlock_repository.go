package store

import (
	"context"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
)

type UnlockRepository struct{ DB *DB }

func (r UnlockRepository) Grant(ctx context.Context, playerID string, reward progression.ProgressionReward) (bool, error) {
	var inserted bool
	err := r.DB.SQL.QueryRowContext(ctx, `INSERT INTO player_unlocks(player_id,reward_id,reward_type,value,amount)
 VALUES($1,$2,$3,$4,$5)
 ON CONFLICT(player_id,reward_id) DO NOTHING
 RETURNING TRUE`, playerID, reward.ID, reward.Type, reward.Value, reward.Amount).Scan(&inserted)
	if err != nil {
		return false, nil
	}
	return inserted, nil
}
