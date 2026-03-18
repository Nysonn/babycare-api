-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListMessagesByConversation :many
SELECT * FROM messages
WHERE conversation_id = $1
ORDER BY sent_at ASC;

-- name: MarkMessagesAsRead :exec
UPDATE messages
SET is_read = TRUE
WHERE conversation_id = $1
  AND sender_id != $2
  AND is_read = FALSE;

-- name: CountUnreadMessagesForUser :one
SELECT COUNT(*) FROM messages m
JOIN conversations c ON c.id = m.conversation_id
WHERE (c.parent_id = $1 OR c.babysitter_id = $1)
  AND m.sender_id != $1
  AND m.is_read = FALSE;

-- name: GetMessageCountByUserPair :one
SELECT COUNT(*) FROM messages m
JOIN conversations c ON c.id = m.conversation_id
WHERE c.parent_id = $1 AND c.babysitter_id = $2
  AND m.sent_at >= NOW() - INTERVAL '30 days';

-- name: GetUserMessageCount :one
SELECT COUNT(*) FROM messages m
JOIN conversations c ON c.id = m.conversation_id
WHERE (c.parent_id = $1 OR c.babysitter_id = $1)
  AND m.sent_at >= NOW() - INTERVAL '30 days';
