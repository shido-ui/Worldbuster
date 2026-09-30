ALTER TABLE economy_accounts ADD COLUMN IF NOT EXISTS owner_account_id UUID UNIQUE REFERENCES accounts(id) ON DELETE CASCADE;
ALTER TABLE economy_accounts ALTER COLUMN id SET DEFAULT gen_random_uuid();

CREATE TABLE IF NOT EXISTS inventory_stacks (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 item_id TEXT NOT NULL,
 quantity INTEGER NOT NULL CHECK (quantity > 0),
 PRIMARY KEY(player_id,item_id)
);