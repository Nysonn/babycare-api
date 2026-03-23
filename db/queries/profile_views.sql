-- name: RecordProfileView :one
INSERT INTO profile_views (babysitter_id, parent_id)
VALUES ($1, $2)
RETURNING id, babysitter_id, parent_id, viewed_at;

-- name: ListProfileViewsForBabysitter :many
SELECT
    pv.id,
    pv.parent_id,
    u.full_name AS parent_name,
    u.email,
    u.phone,
    pp.occupation,
    pp.primary_location,
    pp.preferred_hours,
    pv.viewed_at,
    EXISTS (
        SELECT 1 FROM conversations c
        JOIN messages m ON m.conversation_id = c.id
        WHERE c.parent_id = pv.parent_id
          AND c.babysitter_id = $1
    ) AS has_messaged
FROM profile_views pv
JOIN users u ON u.id = pv.parent_id
LEFT JOIN parent_profiles pp ON pp.user_id = pv.parent_id
WHERE pv.babysitter_id = $1
ORDER BY pv.viewed_at DESC;

-- name: GetWeeklyViewCount :one
SELECT COUNT(*) FROM profile_views
WHERE babysitter_id = $1
  AND viewed_at >= NOW() - INTERVAL '7 days';
