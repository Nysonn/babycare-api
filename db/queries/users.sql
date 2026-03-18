-- name: CreateUser :one
INSERT INTO users (full_name, email, phone, role, password_hash, clerk_user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByClerkID :one
SELECT * FROM users
WHERE clerk_user_id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE role != 'admin' AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE role = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateUserStatus :one
UPDATE users
SET status = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW(), status = 'deleted', updated_at = NOW()
WHERE id = $1;

-- name: CreateAdminUser :one
INSERT INTO users (full_name, email, password_hash, role, status)
VALUES ($1, $2, $3, 'admin', 'active')
RETURNING *;
