CREATE TABLE IF NOT EXISTS achievement_definitions (
 id TEXT PRIMARY KEY,
 title TEXT NOT NULL,
 description TEXT NOT NULL,
 category TEXT NOT NULL,
 points INTEGER NOT NULL DEFAULT 10 CHECK(points>=0),
 active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS player_achievements (
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 achievement_id TEXT NOT NULL REFERENCES achievement_definitions(id),
 progress BIGINT NOT NULL DEFAULT 0 CHECK(progress>=0),
 completed_at TIMESTAMPTZ,
 PRIMARY KEY(player_id,achievement_id)
);
CREATE INDEX IF NOT EXISTS idx_player_achievements_completed ON player_achievements(player_id,completed_at);

CREATE TABLE IF NOT EXISTS ranking_definitions (
 id TEXT PRIMARY KEY,
 title TEXT NOT NULL,
 metric TEXT NOT NULL,
 active BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS ranking_snapshots (
 ranking_id TEXT NOT NULL REFERENCES ranking_definitions(id),
 player_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
 rank INTEGER NOT NULL CHECK(rank>0),
 score BIGINT NOT NULL,
 captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 PRIMARY KEY(ranking_id,player_id,captured_at)
);
CREATE INDEX IF NOT EXISTS idx_ranking_snapshots_lookup ON ranking_snapshots(ranking_id,captured_at DESC,rank);
INSERT INTO achievement_definitions(id,title,description,category,points) VALUES
('first-level','First Level','Reach level 2.','PROGRESSION',10),
('market-participant','Market Participant','Complete a market trade.','ECONOMY',15)
ON CONFLICT(id) DO NOTHING;
INSERT INTO ranking_definitions(id,title,metric) VALUES
('level','Level Ranking','LEVEL'),
('wealth','Wealth Ranking','CASH')
ON CONFLICT(id) DO NOTHING;
