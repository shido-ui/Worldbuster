CREATE TABLE IF NOT EXISTS player_reputation (
 player_id UUID PRIMARY KEY REFERENCES player_profiles(id) ON DELETE CASCADE,
 public_score INTEGER NOT NULL DEFAULT 0 CHECK(public_score BETWEEN -1000 AND 1000),
 trust_score INTEGER NOT NULL DEFAULT 0 CHECK(trust_score BETWEEN -1000 AND 1000),
 notoriety_score INTEGER NOT NULL DEFAULT 0 CHECK(notoriety_score BETWEEN 0 AND 1000),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS reputation_history (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 source TEXT NOT NULL,
 public_delta INTEGER NOT NULL DEFAULT 0,
 trust_delta INTEGER NOT NULL DEFAULT 0,
 notoriety_delta INTEGER NOT NULL DEFAULT 0,
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_reputation_history_player_time ON reputation_history(player_id,created_at DESC);
