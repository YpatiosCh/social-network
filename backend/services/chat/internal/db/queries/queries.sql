-- name: GetConversationMembers :many
SELECT cm2.user_id
FROM conversation_members cm1
JOIN conversation_members cm2
  ON cm2.conversation_id = cm1.conversation_id
WHERE cm1.user_id = $2
  AND cm2.conversation_id = $1
  AND cm2.user_id <> $2
  AND cm2.deleted_at IS NULL;



-- name: GetMessages :many
SELECT m.*
FROM messages m
JOIN conversation_members cm 
  ON cm.conversation_id = m.conversation_id
WHERE m.conversation_id = $1
  AND cm.user_id = $2
  AND m.deleted_at IS NULL
ORDER BY m.created_at ASC
LIMIT $3 OFFSET $4;


-- name: CreateConversation :one
INSERT INTO conversations (group_id)
VALUES ($1)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (conversation_id, sender_id, message_text)
SELECT $1, $2, $3
FROM conversation_members
WHERE conversation_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING *;


-- name: UpdateLastReadMessage :one
UPDATE conversation_members cm
SET last_read_message_id = $3
WHERE cm.conversation_id = $1
  AND cm.user_id = $2
  AND cm.deleted_at IS NULL
RETURNING *;

-- name: AddConversationMember :one
INSERT INTO conversation_members (conversation_id, user_id)
VALUES ($1, $2)
RETURNING *;


-- name: SoftDeleteConversationMember :one
UPDATE conversation_members cm_target
SET deleted_at = NOW()
FROM conversation_members cm_actor
WHERE cm_target.conversation_id = $1
  AND cm_target.user_id = $2
  AND cm_target.deleted_at IS NULL
  AND cm_actor.conversation_id = $1
  AND cm_actor.user_id = $3
  AND cm_actor.deleted_at IS NULL
RETURNING cm_target.*;

-- name: GetUserConversations :many
SELECT 
    c.id AS conversation_id,
    c.group_id,
    c.created_at,
    c.updated_at,
    cm2.user_id AS member_id
FROM conversations c
JOIN conversation_members cm1
    ON cm1.conversation_id = c.id
    AND cm1.user_id = $1
    AND cm1.deleted_at IS NULL
JOIN conversation_members cm2
    ON cm2.conversation_id = c.id
    AND cm2.user_id <> $1
    AND cm2.deleted_at IS NULL
WHERE c.deleted_at IS NULL
AND (
    ($2 IS NULL AND c.group_id IS NULL)
    OR ($2 IS NOT NULL AND c.group_id = $2)
)
ORDER BY c.updated_at DESC
LIMIT $3 OFFSET $4;





