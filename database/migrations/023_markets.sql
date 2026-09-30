CREATE TABLE IF NOT EXISTS market_assets (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 symbol TEXT NOT NULL UNIQUE,
 name TEXT NOT NULL,
 category TEXT NOT NULL,
 base_price BIGINT NOT NULL CHECK(base_price>0),
 current_price BIGINT NOT NULL CHECK(current_price>0),
 supply BIGINT NOT NULL DEFAULT 0 CHECK(supply>=0),
 demand BIGINT NOT NULL DEFAULT 0 CHECK(demand>=0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS market_orders (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 asset_id UUID NOT NULL REFERENCES market_assets(id),
 seller_account_id UUID NOT NULL REFERENCES accounts(id),
 quantity BIGINT NOT NULL CHECK(quantity>0),
 unit_price BIGINT NOT NULL CHECK(unit_price>0),
 status TEXT NOT NULL DEFAULT 'OPEN' CHECK(status IN ('OPEN','FILLED','CANCELLED')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_market_orders_open ON market_orders(asset_id,status,unit_price,created_at);
CREATE TABLE IF NOT EXISTS market_trades (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 asset_id UUID NOT NULL REFERENCES market_assets(id),
 buyer_account_id UUID NOT NULL REFERENCES accounts(id),
 seller_account_id UUID NOT NULL REFERENCES accounts(id),
 quantity BIGINT NOT NULL CHECK(quantity>0),
 unit_price BIGINT NOT NULL CHECK(unit_price>0),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO market_assets(symbol,name,category,base_price) VALUES('WBX-CREDIT','World Credit','currency',100)
ON CONFLICT(symbol) DO NOTHING;
