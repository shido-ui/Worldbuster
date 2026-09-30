CREATE TABLE IF NOT EXISTS player_unlocks (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 reward_id TEXT NOT NULL,
 reward_type TEXT NOT NULL,
 value TEXT NOT NULL DEFAULT '',
 amount INTEGER NOT NULL DEFAULT 0,
 unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY (player_id, reward_id)
);

CREATE INDEX IF NOT EXISTS idx_player_unlocks_reward ON player_unlocks(reward_id);
