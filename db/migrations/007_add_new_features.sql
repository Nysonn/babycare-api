-- +goose Up

-- Feature 1: primary_location for parents
ALTER TABLE parent_profiles ADD COLUMN IF NOT EXISTS primary_location TEXT;

-- Feature 2: gender, availability, currency for babysitters
ALTER TABLE babysitter_profiles ADD COLUMN IF NOT EXISTS gender VARCHAR(10);
ALTER TABLE babysitter_profiles ADD COLUMN IF NOT EXISTS availability TEXT[];
ALTER TABLE babysitter_profiles ADD COLUMN IF NOT EXISTS currency VARCHAR(10) NOT NULL DEFAULT 'UGX';

-- Feature 4: work status (babysitter toggles availability)
ALTER TABLE babysitter_profiles ADD COLUMN IF NOT EXISTS is_available BOOLEAN NOT NULL DEFAULT TRUE;

-- Feature 3: saved babysitters (parents bookmark babysitters)
CREATE TABLE IF NOT EXISTS saved_babysitters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    babysitter_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    saved_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (parent_id, babysitter_id)
);

CREATE INDEX IF NOT EXISTS idx_saved_babysitters_parent_id ON saved_babysitters(parent_id);

-- +goose Down
DROP TABLE IF EXISTS saved_babysitters;
ALTER TABLE babysitter_profiles DROP COLUMN IF EXISTS is_available;
ALTER TABLE babysitter_profiles DROP COLUMN IF EXISTS currency;
ALTER TABLE babysitter_profiles DROP COLUMN IF EXISTS availability;
ALTER TABLE babysitter_profiles DROP COLUMN IF EXISTS gender;
ALTER TABLE parent_profiles DROP COLUMN IF EXISTS primary_location;
