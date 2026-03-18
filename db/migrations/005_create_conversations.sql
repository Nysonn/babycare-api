-- +goose Up
CREATE TABLE conversations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    babysitter_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stream_channel_id TEXT,
    is_locked       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(parent_id, babysitter_id)
);

CREATE INDEX idx_conversations_parent_id ON conversations(parent_id);
CREATE INDEX idx_conversations_babysitter_id ON conversations(babysitter_id);

-- +goose Down
DROP TABLE IF EXISTS conversations;
