-- name: CreateParentProfile :one
INSERT INTO parent_profiles (user_id, location, occupation, preferred_hours, primary_location)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, location, occupation, preferred_hours, created_at, updated_at, primary_location;

-- name: GetParentProfileByUserID :one
SELECT id, user_id, location, occupation, preferred_hours, created_at, updated_at, primary_location
FROM parent_profiles
WHERE user_id = $1;

-- name: UpdateParentProfile :one
UPDATE parent_profiles
SET location = $2, occupation = $3, preferred_hours = $4, primary_location = $5, updated_at = NOW()
WHERE user_id = $1
RETURNING id, user_id, location, occupation, preferred_hours, created_at, updated_at, primary_location;
