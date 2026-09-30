CREATE TABLE IF NOT EXISTS world_event_definitions (
 code TEXT PRIMARY KEY,
 title TEXT NOT NULL,
 description TEXT NOT NULL,
 event_type TEXT NOT NULL,
 cooldown_seconds INTEGER NOT NULL DEFAULT 300 CHECK(cooldown_seconds>=0),
 active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS world_event_instances (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 code TEXT NOT NULL REFERENCES world_event_definitions(code),
 status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK(status IN ('ACTIVE','RESOLVED','EXPIRED')),
 severity INTEGER NOT NULL DEFAULT 1 CHECK(severity BETWEEN 1 AND 5),
 location_id TEXT,
 payload JSONB NOT NULL DEFAULT '{}'::jsonb,
 started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_world_event_instances_status ON world_event_instances(status,started_at DESC);

INSERT INTO world_event_definitions(code,title,description,event_type,cooldown_seconds)
VALUES
('market-fluctuation','Market Fluctuation','A change in local market conditions affects commerce.','ECONOMY',300),
('district-activity','District Activity','A burst of activity changes local world conditions.','SOCIAL',300)
ON CONFLICT(code) DO NOTHING;
