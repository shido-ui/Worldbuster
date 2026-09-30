CREATE TABLE IF NOT EXISTS player_skills (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 skill_id TEXT NOT NULL,
 level INTEGER NOT NULL DEFAULT 1 CHECK (level >= 1),
 xp BIGINT NOT NULL DEFAULT 0 CHECK (xp >= 0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY (player_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_player_skills_skill ON player_skills(skill_id);
CREATE INDEX IF NOT EXISTS idx_player_skills_level ON player_skills(level);
