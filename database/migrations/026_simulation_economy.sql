CREATE TABLE IF NOT EXISTS simulated_economy (
 character_id UUID PRIMARY KEY REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 balance BIGINT NOT NULL DEFAULT 100 CHECK (balance >= 0),
 lifetime_income BIGINT NOT NULL DEFAULT 0 CHECK (lifetime_income >= 0),
 lifetime_spending BIGINT NOT NULL DEFAULT 0 CHECK (lifetime_spending >= 0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS simulated_economy_ledger (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 amount BIGINT NOT NULL CHECK (amount <> 0),
 balance_after BIGINT NOT NULL CHECK (balance_after >= 0),
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_economy_ledger_character_time ON simulated_economy_ledger(character_id,created_at DESC);
