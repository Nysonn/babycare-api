-- name: RecordProfileView :one
INSERT INTO profile_views (babysitter_id, parent_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListProfileViewsForBabysitter :many
SELECT pv.id, pv.parent_id, u.full_name AS parent_name, pv.viewed_at
FROM profile_views pv
JOIN users u ON u.id = pv.parent_id
WHERE pv.babysitter_id = $1
ORDER BY pv.viewed_at DESC;
