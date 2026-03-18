-- +goose Up
CREATE TYPE rate_type AS ENUM ('hourly', 'daily', 'weekly', 'monthly');

CREATE TABLE babysitter_profiles (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location                TEXT,
    national_id_url         TEXT,
    lci_letter_url          TEXT,
    cv_url                  TEXT,
    profile_picture_url     TEXT,
    languages               TEXT[],
    days_per_week           INTEGER,
    hours_per_day           INTEGER,
    rate_type               rate_type,
    rate_amount             NUMERIC(10, 2),
    payment_method          VARCHAR(100),
    is_approved             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_babysitter_profiles_user_id ON babysitter_profiles(user_id);
CREATE INDEX idx_babysitter_profiles_is_approved ON babysitter_profiles(is_approved);

-- +goose Down
DROP TABLE IF EXISTS babysitter_profiles;
DROP TYPE IF EXISTS rate_type;
