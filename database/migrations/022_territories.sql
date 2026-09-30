CREATE TABLE IF NOT EXISTS territories (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 location_id TEXT NOT NULL UNIQUE,
 name TEXT NOT NULL,
 controlling_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
 influence INTEGER NOT NULL DEFAULT 0 CHECK(influence BETWEEN 0 AND 1000),
 stability INTEGER NOT NULL DEFAULT 100 CHECK(stability BETWEEN 0 AND 100),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS territory_influence (
 territory_id UUID NOT NULL REFERENCES territories(id) ON DELETE CASCADE,
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 influence INTEGER NOT NULL DEFAULT 0 CHECK(influence BETWEEN 0 AND 1000),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(territory_id,organization_id)
);
CREATE TABLE IF NOT EXISTS territory_history (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 territory_id UUID NOT NULL REFERENCES territories(id) ON DELETE CASCADE,
 organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
 actor_id UUID,
 delta INTEGER NOT NULL,
 reason TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_territory_history_time ON territory_history(territory_id,created_at DESC);
INSERT INTO territories(location_id,name) VALUES
 ('central','Central District')
ON CONFLICT(location_id) DO NOTHING;
