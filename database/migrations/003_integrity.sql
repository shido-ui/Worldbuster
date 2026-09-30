CREATE UNIQUE INDEX IF NOT EXISTS uq_accounts_username_ci ON accounts (LOWER(username));
ALTER TABLE player_profiles DROP CONSTRAINT IF EXISTS player_profiles_energy_valid;
ALTER TABLE player_profiles ADD CONSTRAINT player_profiles_energy_valid CHECK (energy BETWEEN 0 AND 100);
CREATE INDEX IF NOT EXISTS idx_characters_controller_type ON characters(controller_type);
