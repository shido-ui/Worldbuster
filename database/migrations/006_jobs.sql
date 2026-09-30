-- Legacy job schema retained for historical compatibility.
-- Persistent player employment uses the authoritative UUID job schema introduced later.
CREATE TABLE IF NOT EXISTS legacy_jobs (id TEXT PRIMARY KEY,name TEXT NOT NULL,department TEXT NOT NULL,description TEXT NOT NULL DEFAULT '');
CREATE TABLE IF NOT EXISTS legacy_job_positions (id TEXT PRIMARY KEY,job_id TEXT NOT NULL REFERENCES legacy_jobs(id) ON DELETE CASCADE,name TEXT NOT NULL,level INT NOT NULL CHECK(level>0),base_salary BIGINT NOT NULL CHECK(base_salary>=0),required_education INT NOT NULL DEFAULT 0,required_stat INT NOT NULL DEFAULT 0);
CREATE TABLE IF NOT EXISTS legacy_employments (character_id UUID PRIMARY KEY,job_id TEXT NOT NULL REFERENCES legacy_jobs(id),position_id TEXT NOT NULL REFERENCES legacy_job_positions(id),experience INT NOT NULL DEFAULT 0 CHECK(experience>=0),started_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
