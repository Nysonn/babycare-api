-- +goose Up
CREATE TABLE profile_views (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    babysitter_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    viewed_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profile_views_babysitter_id ON profile_views(babysitter_id);
CREATE INDEX idx_profile_views_parent_id ON profile_views(parent_id);

-- +goose Down
DROP TABLE IF EXISTS profile_views;
