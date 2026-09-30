ALTER TABLE territory_influence ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE TABLE IF NOT EXISTS territory_effects (
 territory_id UUID PRIMARY KEY REFERENCES territories(id) ON DELETE CASCADE,
 commerce_modifier_bps INTEGER NOT NULL DEFAULT 0 CHECK(commerce_modifier_bps BETWEEN -5000 AND 5000),
 safety_modifier INTEGER NOT NULL DEFAULT 0 CHECK(safety_modifier BETWEEN -100 AND 100),
 stability_modifier INTEGER NOT NULL DEFAULT 0 CHECK(stability_modifier BETWEEN -100 AND 100),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS territory_consequence_cycles (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 territories_processed INTEGER NOT NULL,
 control_changes INTEGER NOT NULL,
 stability_changes INTEGER NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_territory_influence_activity ON territory_influence(territory_id,last_activity_at);
