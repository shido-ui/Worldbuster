CREATE TABLE IF NOT EXISTS admin_roles (
 account_id UUID PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
 role TEXT NOT NULL CHECK(role IN ('WORLD_ADMIN','WORLD_OPERATOR')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS world_control_audit (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 account_id UUID REFERENCES accounts(id),
 action TEXT NOT NULL,
 target_type TEXT NOT NULL,
 target_id TEXT,
 payload JSONB NOT NULL DEFAULT '{}'::jsonb,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_world_control_audit_time ON world_control_audit(created_at DESC);

CREATE TABLE IF NOT EXISTS world_control_flags (
 key TEXT PRIMARY KEY,
 value JSONB NOT NULL,
 updated_by UUID REFERENCES accounts(id),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
