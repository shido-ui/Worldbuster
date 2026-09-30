CREATE TABLE IF NOT EXISTS simulated_character_memory (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 memory_type TEXT NOT NULL,
 target_id TEXT,
 event TEXT NOT NULL,
 importance INTEGER NOT NULL DEFAULT 50 CHECK(importance BETWEEN 0 AND 100),
 sentiment INTEGER NOT NULL DEFAULT 0 CHECK(sentiment BETWEEN -100 AND 100),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_memory_character_time ON simulated_character_memory(character_id,created_at DESC);
CREATE TABLE IF NOT EXISTS simulated_character_actions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 action_type TEXT NOT NULL,
 target_id TEXT,
 decision_score INTEGER NOT NULL DEFAULT 0,
 executed BOOLEAN NOT NULL DEFAULT FALSE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_actions_character_time ON simulated_character_actions(character_id,created_at DESC);
