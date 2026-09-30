ALTER TABLE achievement_definitions ADD COLUMN IF NOT EXISTS target_progress BIGINT NOT NULL DEFAULT 1 CHECK(target_progress>0);

UPDATE achievement_definitions SET target_progress=1 WHERE id IN ('first-level','market-participant');

CREATE TABLE IF NOT EXISTS market_asset_inventory (
 account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
 asset_id UUID NOT NULL REFERENCES market_assets(id) ON DELETE CASCADE,
 quantity BIGINT NOT NULL DEFAULT 0 CHECK(quantity>=0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(account_id,asset_id)
);
CREATE INDEX IF NOT EXISTS idx_market_asset_inventory_asset ON market_asset_inventory(asset_id);

INSERT INTO market_asset_inventory(account_id,asset_id,quantity)
SELECT a.id,m.id,0 FROM accounts a CROSS JOIN market_assets m
ON CONFLICT DO NOTHING;
