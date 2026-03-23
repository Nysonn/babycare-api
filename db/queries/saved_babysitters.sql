-- name: SaveBabysitter :exec
INSERT INTO saved_babysitters (parent_id, babysitter_id)
VALUES ($1, $2)
ON CONFLICT (parent_id, babysitter_id) DO NOTHING;

-- name: UnsaveBabysitter :exec
DELETE FROM saved_babysitters
WHERE parent_id = $1 AND babysitter_id = $2;

-- name: ListSavedBabysitters :many
SELECT u.id, u.full_name, u.email, u.phone, u.status,
       bp.location, bp.profile_picture_url, bp.languages,
       bp.days_per_week, bp.hours_per_day, bp.rate_type,
       bp.rate_amount, bp.payment_method, bp.is_approved,
       bp.gender, bp.availability, bp.currency, bp.is_available
FROM saved_babysitters sb
JOIN users u ON u.id = sb.babysitter_id
JOIN babysitter_profiles bp ON u.id = bp.user_id
WHERE sb.parent_id = $1
  AND u.deleted_at IS NULL
  AND u.status = 'active'
ORDER BY sb.saved_at DESC;
