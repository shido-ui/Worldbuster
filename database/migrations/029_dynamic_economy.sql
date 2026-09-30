CREATE TABLE IF NOT EXISTS economy_price_history (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 asset_id UUID NOT NULL REFERENCES market_assets(id) ON DELETE CASCADE,
 old_price BIGINT NOT NULL CHECK(old_price>0),
 new_price BIGINT NOT NULL CHECK(new_price>0),
 supply BIGINT NOT NULL CHECK(supply>=0),
 demand BIGINT NOT NULL CHECK(demand>=0),
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_economy_price_history_asset_time ON economy_price_history(asset_id,created_at DESC);

CREATE TABLE IF NOT EXISTS economy_balance_cycles (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 processed_assets INTEGER NOT NULL DEFAULT 0,
 price_changes INTEGER NOT NULL DEFAULT 0,
 total_supply BIGINT NOT NULL DEFAULT 0,
 total_demand BIGINT NOT NULL DEFAULT 0,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE market_assets ADD COLUMN IF NOT EXISTS volatility_bps INTEGER NOT NULL DEFAULT 500 CHECK(volatility_bps BETWEEN 0 AND 5000);
