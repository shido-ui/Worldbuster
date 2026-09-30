ALTER TABLE player_employment ADD COLUMN IF NOT EXISTS last_paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
CREATE INDEX IF NOT EXISTS idx_player_employment_salary_due ON player_employment(last_paid_at);