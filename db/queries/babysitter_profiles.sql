-- name: CreateBabysitterProfile :one
INSERT INTO babysitter_profiles (user_id, location, national_id_url, lci_letter_url, cv_url, profile_picture_url, languages, days_per_week, hours_per_day, rate_type, rate_amount, payment_method)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetBabysitterProfileByUserID :one
SELECT * FROM babysitter_profiles
WHERE user_id = $1;

-- name: UpdateBabysitterProfile :one
UPDATE babysitter_profiles
SET location = $2, languages = $3, days_per_week = $4, hours_per_day = $5,
    rate_type = $6, rate_amount = $7, payment_method = $8,
    profile_picture_url = $9, updated_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: ApproveBabysitter :one
UPDATE babysitter_profiles
SET is_approved = TRUE, updated_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: ListApprovedBabysitters :many
SELECT u.id, u.full_name, u.email, u.phone, u.status,
       bp.location, bp.profile_picture_url, bp.languages,
       bp.days_per_week, bp.hours_per_day, bp.rate_type,
       bp.rate_amount, bp.payment_method, bp.is_approved
FROM users u
JOIN babysitter_profiles bp ON u.id = bp.user_id
WHERE u.role = 'babysitter'
  AND u.deleted_at IS NULL
  AND u.status = 'active'
  AND bp.is_approved = TRUE
ORDER BY u.created_at DESC;
