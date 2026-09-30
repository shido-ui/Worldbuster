CREATE TABLE IF NOT EXISTS simulated_businesses (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 owner_character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 name TEXT NOT NULL,
 business_type TEXT NOT NULL,
 cash_balance BIGINT NOT NULL DEFAULT 0 CHECK(cash_balance>=0),
 level INTEGER NOT NULL DEFAULT 1 CHECK(level>=1),
 reputation INTEGER NOT NULL DEFAULT 0 CHECK(reputation BETWEEN -1000 AND 1000),
 employees INTEGER NOT NULL DEFAULT 0 CHECK(employees>=0),
 revenue_total BIGINT NOT NULL DEFAULT 0 CHECK(revenue_total>=0),
 expense_total BIGINT NOT NULL DEFAULT 0 CHECK(expense_total>=0),
 active BOOLEAN NOT NULL DEFAULT TRUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sim_business_owner ON simulated_businesses(owner_character_id);
CREATE INDEX IF NOT EXISTS idx_sim_business_type_active ON simulated_businesses(business_type,active);

CREATE TABLE IF NOT EXISTS simulated_business_ledger (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 business_id UUID NOT NULL REFERENCES simulated_businesses(id) ON DELETE CASCADE,
 amount BIGINT NOT NULL CHECK(amount<>0),
 balance_after BIGINT NOT NULL CHECK(balance_after>=0),
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_business_ledger_time ON simulated_business_ledger(business_id,created_at DESC);

CREATE TABLE IF NOT EXISTS simulated_business_employment (
 business_id UUID NOT NULL REFERENCES simulated_businesses(id) ON DELETE CASCADE,
 character_id UUID NOT NULL REFERENCES simulated_characters(character_id) ON DELETE CASCADE,
 role TEXT NOT NULL,
 wage BIGINT NOT NULL CHECK(wage>=0),
 active BOOLEAN NOT NULL DEFAULT TRUE,
 hired_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(business_id,character_id)
);
CREATE INDEX IF NOT EXISTS idx_sim_business_employees ON simulated_business_employment(character_id,active);
