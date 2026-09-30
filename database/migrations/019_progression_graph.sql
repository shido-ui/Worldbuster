ALTER TABLE jobs ADD COLUMN IF NOT EXISTS required_education INTEGER NOT NULL DEFAULT 0 CHECK (required_education >= 0);
CREATE INDEX IF NOT EXISTS idx_jobs_progression_requirements ON jobs(required_level,required_education);
UPDATE jobs SET required_education=0 WHERE required_education IS NULL;