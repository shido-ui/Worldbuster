CREATE TABLE IF NOT EXISTS world_events (id UUID PRIMARY KEY,type TEXT NOT NULL,actor_id UUID,target_id UUID,payload JSONB NOT NULL DEFAULT '{}'::jsonb,created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX IF NOT EXISTS idx_world_events_time ON world_events(created_at DESC);
CREATE TABLE IF NOT EXISTS notifications (id UUID PRIMARY KEY,character_id UUID NOT NULL,type TEXT NOT NULL,title TEXT NOT NULL,body TEXT NOT NULL,event_id UUID REFERENCES world_events(id),created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),read_at TIMESTAMPTZ);
CREATE INDEX IF NOT EXISTS idx_notifications_character_time ON notifications(character_id,created_at DESC);
