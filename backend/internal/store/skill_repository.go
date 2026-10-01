package store

import (
	"context"
	"github.com/shido-ui/Worldbuster/backend/internal/progression"
)

type SkillRepository struct{ DB *DB }

func (r SkillRepository) Upsert(ctx context.Context, playerID string, state progression.SkillState) (progression.SkillState, error) {
	var out progression.SkillState
	err := r.DB.SQL.QueryRowContext(ctx, `INSERT INTO player_skills(player_id,skill_id,level,xp)
 VALUES($1,$2,$3,$4)
 ON CONFLICT(player_id,skill_id) DO UPDATE SET level=EXCLUDED.level,xp=EXCLUDED.xp,updated_at=NOW()
 RETURNING skill_id,level,xp`, playerID, state.SkillID, state.Level, state.XP).Scan(&out.SkillID, &out.Level, &out.XP)
	return out, err
}

func (r SkillRepository) Get(ctx context.Context, playerID, skillID string) (progression.SkillState, error) {
	var out progression.SkillState
	err := r.DB.SQL.QueryRowContext(ctx, `SELECT skill_id,level,xp FROM player_skills WHERE player_id=$1 AND skill_id=$2`, playerID, skillID).Scan(&out.SkillID, &out.Level, &out.XP)
	return out, err
}
