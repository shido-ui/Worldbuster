ALTER TABLE organizations ADD COLUMN IF NOT EXISTS level INTEGER NOT NULL DEFAULT 1 CHECK(level>=1);
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS reputation INTEGER NOT NULL DEFAULT 0 CHECK(reputation BETWEEN -1000 AND 1000);
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS treasury BIGINT NOT NULL DEFAULT 0 CHECK(treasury>=0);
CREATE TABLE IF NOT EXISTS organization_reputation (
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 score INTEGER NOT NULL DEFAULT 0 CHECK(score BETWEEN -1000 AND 1000),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(organization_id,player_id)
);
CREATE TABLE IF NOT EXISTS organization_history (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 actor_id UUID,
 event_type TEXT NOT NULL,
 payload JSONB NOT NULL DEFAULT '{}'::jsonb,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_org_history_time ON organization_history(organization_id,created_at DESC);
