CREATE UNIQUE INDEX IF NOT EXISTS uq_accounts_username_ci ON accounts (LOWER(username));
ALTER TABLE player_profiles DROP CONSTRAINT IF EXISTS player_profiles_energy_valid;
ALTER TABLE player_profiles ADD CONSTRAINT player_profiles_energy_valid CHECK (energy BETWEEN 0 AND 100);
-- The historical characters table is not part of the authoritative migration chain.
-- The obsolete index was removed so fresh databases do not depend on an undeclared table.
