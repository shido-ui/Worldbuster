ALTER TABLE player_profiles ADD COLUMN IF NOT EXISTS location_id TEXT NOT NULL DEFAULT 'central';
CREATE INDEX IF NOT EXISTS idx_player_profiles_location ON player_profiles(location_id);
