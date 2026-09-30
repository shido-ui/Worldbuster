CREATE TABLE IF NOT EXISTS production_recipes (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 name TEXT NOT NULL UNIQUE,
 business_type TEXT NOT NULL,
 min_level INTEGER NOT NULL DEFAULT 1 CHECK(min_level>=1),
 input_item TEXT NOT NULL,
 input_quantity INTEGER NOT NULL CHECK(input_quantity>0),
 output_item TEXT NOT NULL,
 output_quantity INTEGER NOT NULL CHECK(output_quantity>0),
 labor_cost BIGINT NOT NULL DEFAULT 1 CHECK(labor_cost>=0),
 active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS simulated_business_inventory (
 business_id UUID NOT NULL REFERENCES simulated_businesses(id) ON DELETE CASCADE,
 item_id TEXT NOT NULL,
 quantity BIGINT NOT NULL DEFAULT 0 CHECK(quantity>=0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(business_id,item_id)
);

CREATE TABLE IF NOT EXISTS simulated_production_runs (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 business_id UUID NOT NULL REFERENCES simulated_businesses(id) ON DELETE CASCADE,
 recipe_id UUID NOT NULL REFERENCES production_recipes(id),
 input_item TEXT NOT NULL,
 input_quantity INTEGER NOT NULL,
 output_item TEXT NOT NULL,
 output_quantity INTEGER NOT NULL,
 labor_cost BIGINT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sim_production_runs_time ON simulated_production_runs(created_at DESC);

CREATE TABLE IF NOT EXISTS world_supply_signals (
 item_id TEXT PRIMARY KEY,
 supply BIGINT NOT NULL DEFAULT 0 CHECK(supply>=0),
 demand BIGINT NOT NULL DEFAULT 0 CHECK(demand>=0),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO production_recipes(name,business_type,min_level,input_item,input_quantity,output_item,output_quantity,labor_cost)
VALUES
 ('Water Bottling','PRODUCTION',1,'raw-water',5,'bottled-water',5,2),
 ('Basic Textile Production','PRODUCTION',1,'fiber',4,'fabric',2,3),
 ('Tool Assembly','PRODUCTION',2,'fabric',2,'basic-tools',1,6)
ON CONFLICT(name) DO NOTHING;

INSERT INTO world_supply_signals(item_id,supply,demand)
VALUES ('raw-water',100,80),('bottled-water',20,80),('fiber',80,50),('fabric',15,55),('basic-tools',5,35)
ON CONFLICT(item_id) DO NOTHING;
