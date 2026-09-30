CREATE TABLE IF NOT EXISTS simulated_relationships (
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 target_id UUID NOT NULL,
 affinity INTEGER NOT NULL DEFAULT 0 CHECK (affinity BETWEEN -1000 AND 1000),
 trust INTEGER NOT NULL DEFAULT 0 CHECK (trust BETWEEN -1000 AND 1000),
 familiarity INTEGER NOT NULL DEFAULT 0 CHECK (familiarity BETWEEN 0 AND 1000),
 interactions INTEGER NOT NULL DEFAULT 0 CHECK (interactions >= 0),
 last_interaction TIMESTAMPTZ,
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY (character_id,target_id),
 CHECK (character_id <> target_id)
);
CREATE INDEX IF NOT EXISTS idx_sim_relationships_target ON simulated_relationships(target_id);
CREATE INDEX IF NOT EXISTS idx_sim_relationships_affinity ON simulated_relationships(character_id,affinity DESC);
CREATE TABLE IF NOT EXISTS simulated_social_events (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 target_id UUID NOT NULL,
 event_type TEXT NOT NULL,
 affinity_delta INTEGER NOT NULL DEFAULT 0,
 trust_delta INTEGER NOT NULL DEFAULT 0,
 familiarity_delta INTEGER NOT NULL DEFAULT 0,
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_social_events_character_time ON simulated_social_events(character_id,created_at DESC);
