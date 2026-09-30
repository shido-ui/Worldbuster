CREATE TABLE IF NOT EXISTS relationships (character_id UUID NOT NULL,target_id UUID NOT NULL,kind TEXT NOT NULL,score INT NOT NULL DEFAULT 0,updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),PRIMARY KEY(character_id,target_id),CHECK(character_id<>target_id));
CREATE TABLE IF NOT EXISTS messages (id UUID PRIMARY KEY,from_id UUID NOT NULL,to_id UUID NOT NULL,body TEXT NOT NULL,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),read_at TIMESTAMPTZ);
CREATE INDEX IF NOT EXISTS idx_messages_inbox ON messages(to_id,created_at DESC);
