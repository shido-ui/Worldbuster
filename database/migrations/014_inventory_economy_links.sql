ALTER TABLE economy_accounts ADD COLUMN IF NOT EXISTS owner_account_id UUID UNIQUE REFERENCES accounts(id) ON DELETE CASCADE;
CREATE TABLE IF NOT EXISTS inventory_stacks (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 item_id TEXT NOT NULL,
 quantity INTEGER NOT NULL CHECK (quantity > 0),
 PRIMARY KEY(player_id,item_id)
);
CREATE INDEX IF NOT EXISTS idx_inventory_stacks_player ON inventory_stacks(player_id);
