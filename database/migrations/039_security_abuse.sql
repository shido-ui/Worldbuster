CREATE TABLE IF NOT EXISTS request_rate_buckets (
 bucket_key TEXT PRIMARY KEY,
 window_started_at TIMESTAMPTZ NOT NULL,
 request_count INTEGER NOT NULL DEFAULT 0 CHECK(request_count>=0)
);
CREATE TABLE IF NOT EXISTS abuse_audit (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
 category TEXT NOT NULL,
 action TEXT NOT NULL,
 detail JSONB NOT NULL DEFAULT '{}'::jsonb,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_abuse_audit_account_time ON abuse_audit(account_id,created_at DESC);
