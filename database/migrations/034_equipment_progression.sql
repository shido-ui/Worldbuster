CREATE TABLE IF NOT EXISTS item_definitions (
 id TEXT PRIMARY KEY,
 name TEXT NOT NULL,
 category TEXT NOT NULL,
 stackable BOOLEAN NOT NULL DEFAULT TRUE,
 max_stack INTEGER NOT NULL DEFAULT 99 CHECK(max_stack>0),
 slot TEXT,
 power_bonus INTEGER NOT NULL DEFAULT 0 CHECK(power_bonus>=0),
 defense_bonus INTEGER NOT NULL DEFAULT 0 CHECK(defense_bonus>=0),
 speed_bonus INTEGER NOT NULL DEFAULT 0 CHECK(speed_bonus>=0),
 durability_max INTEGER NOT NULL DEFAULT 0 CHECK(durability_max>=0),
 required_level INTEGER NOT NULL DEFAULT 1 CHECK(required_level>=1),
 active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS player_equipment (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 slot TEXT NOT NULL,
 item_id TEXT NOT NULL REFERENCES item_definitions(id),
 durability INTEGER NOT NULL DEFAULT 0 CHECK(durability>=0),
 equipped_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(player_id,slot)
);

CREATE TABLE IF NOT EXISTS item_progression (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 item_id TEXT NOT NULL REFERENCES item_definitions(id),
 upgrade_level INTEGER NOT NULL DEFAULT 0 CHECK(upgrade_level>=0),
 experience BIGINT NOT NULL DEFAULT 0 CHECK(experience>=0),
 PRIMARY KEY(player_id,item_id)
);

INSERT INTO item_definitions(id,name,category,stackable,max_stack,slot,power_bonus,defense_bonus,speed_bonus,durability_max,required_level)
VALUES
('training-gear','Training Gear','EQUIPMENT',FALSE,1,'PRIMARY',2,0,0,100,1),
('reinforced-vest','Reinforced Vest','EQUIPMENT',FALSE,1,'ARMOR',0,5,0,100,1),
('field-boots','Field Boots','EQUIPMENT',FALSE,1,'FOOTWEAR',0,1,3,80,1)
ON CONFLICT(id) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_item_progression_player ON item_progression(player_id);
