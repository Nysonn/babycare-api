-- +goose Up
ALTER TABLE parent_profiles
ADD COLUMN profile_picture_url TEXT;

-- +goose Down
ALTER TABLE parent_profiles
DROP COLUMN IF EXISTS profile_picture_url;