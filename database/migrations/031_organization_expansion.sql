ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS permissions JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS influence BIGINT NOT NULL DEFAULT 0 CHECK(influence>=0);
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS activity_score INTEGER NOT NULL DEFAULT 0 CHECK(activity_score BETWEEN 0 AND 1000);

CREATE TABLE IF NOT EXISTS organization_alliances (
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 target_organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 status TEXT NOT NULL DEFAULT 'PROPOSED' CHECK(status IN ('PROPOSED','ACTIVE','ENDED')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(organization_id,target_organization_id),
 CHECK(organization_id<>target_organization_id)
);

CREATE TABLE IF NOT EXISTS organization_goals (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 goal_type TEXT NOT NULL,
 target_value BIGINT NOT NULL DEFAULT 1 CHECK(target_value>0),
 progress_value BIGINT NOT NULL DEFAULT 0 CHECK(progress_value>=0),
 status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK(status IN ('ACTIVE','COMPLETED','CANCELLED')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_org_goals_status ON organization_goals(organization_id,status);
CREATE INDEX IF NOT EXISTS idx_org_alliances_target ON organization_alliances(target_organization_id);
