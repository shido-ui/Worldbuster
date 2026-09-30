ALTER TABLE player_profiles ADD COLUMN IF NOT EXISTS education INTEGER NOT NULL DEFAULT 0 CHECK (education >= 0);

CREATE TABLE IF NOT EXISTS education_courses (
 id TEXT PRIMARY KEY,
 name TEXT NOT NULL UNIQUE,
 duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
 education_gain INTEGER NOT NULL CHECK (education_gain > 0),
 required_level INTEGER NOT NULL DEFAULT 1 CHECK (required_level >= 1)
);

CREATE TABLE IF NOT EXISTS player_training (
 player_id UUID PRIMARY KEY REFERENCES player_profiles(id) ON DELETE CASCADE,
 course_id TEXT NOT NULL REFERENCES education_courses(id),
 started_at TIMESTAMPTZ NOT NULL,
 completes_at TIMESTAMPTZ NOT NULL,
 completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_player_training_due ON player_training(completes_at) WHERE completed_at IS NULL;

INSERT INTO education_courses(id,name,duration_minutes,education_gain,required_level)
VALUES ('orientation','World Orientation',10,1,1)
ON CONFLICT(id) DO NOTHING;