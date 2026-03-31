-- name: CreateConversation :one
INSERT INTO conversations (parent_id, babysitter_id, stream_channel_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetConversationByID :one
SELECT * FROM conversations
WHERE id = $1;

-- name: GetConversationByParticipants :one
SELECT * FROM conversations
WHERE parent_id = $1 AND babysitter_id = $2;

-- name: ListConversationsForUser :many
SELECT c.*,
       u.full_name AS other_user_name,
       COALESCE(
           CASE WHEN c.parent_id = $1 THEN bp.profile_picture_url
                ELSE pp.profile_picture_url
           END,
           ''
       ) AS other_user_profile_picture_url
FROM conversations c
JOIN users u ON (
    CASE WHEN c.parent_id = $1 THEN u.id = c.babysitter_id
         ELSE u.id = c.parent_id
    END
)
LEFT JOIN babysitter_profiles bp ON bp.user_id = c.babysitter_id
LEFT JOIN parent_profiles pp ON pp.user_id = c.parent_id
WHERE c.parent_id = $1 OR c.babysitter_id = $1
ORDER BY c.updated_at DESC;

-- name: LockConversation :exec
UPDATE conversations
SET is_locked = TRUE, updated_at = NOW()
WHERE id = $1;

-- name: LockConversationsByUser :exec
UPDATE conversations
SET is_locked = TRUE, updated_at = NOW()
WHERE babysitter_id = $1 OR parent_id = $1;
