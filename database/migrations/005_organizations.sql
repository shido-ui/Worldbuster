CREATE TABLE IF NOT EXISTS organizations (id UUID PRIMARY KEY,name TEXT NOT NULL,type TEXT NOT NULL,owner_id UUID NOT NULL, max_members INT NOT NULL CHECK (max_members > 0),created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE IF NOT EXISTS organization_members (organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,character_id UUID NOT NULL,role TEXT NOT NULL,joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),PRIMARY KEY(organization_id,character_id));
CREATE INDEX IF NOT EXISTS idx_org_members_character ON organization_members(character_id);
