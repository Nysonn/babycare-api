-- name: CreateParentProfile :one
INSERT INTO parent_profiles (user_id, location, occupation, preferred_hours)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetParentProfileByUserID :one
SELECT * FROM parent_profiles
WHERE user_id = $1;

-- name: UpdateParentProfile :one
UPDATE parent_profiles
SET location = $2, occupation = $3, preferred_hours = $4, updated_at = NOW()
WHERE user_id = $1
RETURNING *;
