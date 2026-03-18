-- +goose Up
CREATE TABLE parent_profiles (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location          TEXT,
    occupation        VARCHAR(255),
    preferred_hours   VARCHAR(255),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_parent_profiles_user_id ON parent_profiles(user_id);

-- +goose Down
DROP TABLE IF EXISTS parent_profiles;
