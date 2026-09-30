CREATE TABLE IF NOT EXISTS missions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 code TEXT NOT NULL UNIQUE,
 title TEXT NOT NULL,
 description TEXT NOT NULL,
 mission_type TEXT NOT NULL,
 min_level INTEGER NOT NULL DEFAULT 1 CHECK(min_level>=1),
 reward_cash BIGINT NOT NULL DEFAULT 0 CHECK(reward_cash>=0),
 reward_xp BIGINT NOT NULL DEFAULT 0 CHECK(reward_xp>=0),
 target_type TEXT,
 target_value BIGINT NOT NULL DEFAULT 1 CHECK(target_value>0),
 active BOOLEAN NOT NULL DEFAULT TRUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS player_missions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 mission_id UUID NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK(status IN ('ACTIVE','COMPLETED','ABANDONED')),
 progress BIGINT NOT NULL DEFAULT 0 CHECK(progress>=0),
 accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 completed_at TIMESTAMPTZ,
 UNIQUE(mission_id,player_id)
);

CREATE TABLE IF NOT EXISTS contracts (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 issuer_type TEXT NOT NULL CHECK(issuer_type IN ('PLAYER','ORGANIZATION','WORLD')),
 issuer_id UUID,
 title TEXT NOT NULL,
 description TEXT NOT NULL,
 target_type TEXT NOT NULL,
 target_value BIGINT NOT NULL DEFAULT 1 CHECK(target_value>0),
 reward_cash BIGINT NOT NULL DEFAULT 0 CHECK(reward_cash>=0),
 reward_xp BIGINT NOT NULL DEFAULT 0 CHECK(reward_xp>=0),
 status TEXT NOT NULL DEFAULT 'OPEN' CHECK(status IN ('OPEN','ACCEPTED','COMPLETED','EXPIRED','CANCELLED')),
 expires_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contract_claims (
 contract_id UUID PRIMARY KEY REFERENCES contracts(id) ON DELETE CASCADE,
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_missions_active ON missions(active,min_level);
CREATE INDEX IF NOT EXISTS idx_player_missions_player ON player_missions(player_id,status);
CREATE INDEX IF NOT EXISTS idx_contracts_status_expiry ON contracts(status,expires_at);
INSERT INTO missions(code,title,description,mission_type,min_level,reward_cash,reward_xp,target_type,target_value)
VALUES
('welcome-steps','First Steps','Complete a small amount of productive activity in the world.','TUTORIAL',1,50,100,'activity',1),
('market-observer','Market Observer','Complete a market interaction and observe the changing economy.','ECONOMY',1,75,125,'market_action',1)
ON CONFLICT(code) DO NOTHING;
