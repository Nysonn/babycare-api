-- name: CreateReport :one
INSERT INTO reports (reporter_id, reported_user_id, report_type, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListReports :many
SELECT
    r.*,
    reporter.full_name   AS reporter_name,
    reporter.email       AS reporter_email,
    reported.full_name   AS reported_name,
    reported.email       AS reported_email
FROM reports r
JOIN users reporter ON reporter.id = r.reporter_id
JOIN users reported ON reported.id = r.reported_user_id
ORDER BY r.created_at DESC;

-- name: ListReportsByStatus :many
SELECT
    r.*,
    reporter.full_name   AS reporter_name,
    reporter.email       AS reporter_email,
    reported.full_name   AS reported_name,
    reported.email       AS reported_email
FROM reports r
JOIN users reporter ON reporter.id = r.reporter_id
JOIN users reported ON reported.id = r.reported_user_id
WHERE r.status = $1
ORDER BY r.created_at DESC;

-- name: UpdateReportStatus :one
UPDATE reports
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetReport :one
SELECT
    r.*,
    reporter.full_name   AS reporter_name,
    reporter.email       AS reporter_email,
    reported.full_name   AS reported_name,
    reported.email       AS reported_email
FROM reports r
JOIN users reporter ON reporter.id = r.reporter_id
JOIN users reported ON reported.id = r.reported_user_id
WHERE r.id = $1;
