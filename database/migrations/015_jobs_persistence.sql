-- Reconcile the original text job schema with the authoritative UUID job schema.
DO $$
DECLARE id_type TEXT;
BEGIN
 SELECT data_type INTO id_type
 FROM information_schema.columns
 WHERE table_schema='public' AND table_name='jobs' AND column_name='id';
 IF id_type='text' THEN
  ALTER TABLE IF EXISTS employments RENAME TO legacy_employments;
  ALTER TABLE IF EXISTS job_positions RENAME TO legacy_job_positions;
  ALTER TABLE IF EXISTS jobs RENAME TO legacy_jobs;
 END IF;
END $$;

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
