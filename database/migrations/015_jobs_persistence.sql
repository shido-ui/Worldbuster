CREATE TABLE IF NOT EXISTS jobs (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 name TEXT NOT NULL UNIQUE,
 department TEXT NOT NULL,
 base_salary BIGINT NOT NULL DEFAULT 0 CHECK (base_salary >= 0),
 required_level INTEGER NOT NULL DEFAULT 1 CHECK (required_level >= 1),
 required_stat INTEGER NOT NULL DEFAULT 0 CHECK (required_stat >= 0)
);
CREATE TABLE IF NOT EXISTS player_employment (
 player_id UUID PRIMARY KEY REFERENCES player_profiles(id) ON DELETE CASCADE,
 job_id UUID NOT NULL REFERENCES jobs(id),
 hired_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
