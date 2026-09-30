CREATE TABLE IF NOT EXISTS combat_profiles (
 player_id UUID PRIMARY KEY REFERENCES player_profiles(id) ON DELETE CASCADE,
 rating INTEGER NOT NULL DEFAULT 1000 CHECK(rating>=0),
 wins INTEGER NOT NULL DEFAULT 0 CHECK(wins>=0),
 losses INTEGER NOT NULL DEFAULT 0 CHECK(losses>=0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS combat_matches (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 attacker_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 defender_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 status TEXT NOT NULL DEFAULT 'RESOLVED' CHECK(status IN ('RESOLVED','CANCELLED')),
 winner_id UUID REFERENCES player_profiles(id) ON DELETE SET NULL,
 attacker_score INTEGER NOT NULL CHECK(attacker_score>=0),
 defender_score INTEGER NOT NULL CHECK(defender_score>=0),
 energy_cost INTEGER NOT NULL CHECK(energy_cost>=0),
 xp_reward BIGINT NOT NULL DEFAULT 0 CHECK(xp_reward>=0),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 CHECK(attacker_id<>defender_id)
);

CREATE TABLE IF NOT EXISTS combat_cooldowns (
 player_id UUID PRIMARY KEY REFERENCES player_profiles(id) ON DELETE CASCADE,
 available_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_combat_matches_player ON combat_matches(attacker_id,created_at DESC);
CREATE INDEX IF NOT EXISTS idx_combat_matches_defender ON combat_matches(defender_id,created_at DESC);
