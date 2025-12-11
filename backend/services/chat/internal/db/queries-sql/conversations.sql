-- name: CreatePrivateConv :one
-- Creates a Conversation if and only if a conversation between the same 2 users does not exist.
-- Returns NULL if a duplicate DM exists (sqlc will error if RETURNING finds no rows).
WITH existing AS (
    SELECT c.id
    FROM conversations c
    JOIN conversation_members cm1 ON cm1.conversation_id = c.id AND cm1.user_id = $1
    JOIN conversation_members cm2 ON cm2.conversation_id = c.id AND cm2.user_id = $2
    WHERE c.group_id IS NULL
)
INSERT INTO conversations (group_id)
SELECT NULL
WHERE NOT EXISTS (SELECT 1 FROM existing)
RETURNING id;


-- name: CreateGroupConv :one
INSERT INTO conversations (group_id)
VALUES ($1)
RETURNING id;



-- name: GetUserPrivateConversations :many
-- param: UserID int8
-- param: GroupID int8?
-- param: Limit int4
-- param: Offset int4
SELECT 
    c.id AS conversation_id,
    c.group_id,
    c.created_at,
    c.updated_at
FROM conversations c
JOIN conversation_members cm
    ON cm.conversation_id = c.id
    AND cm.user_id = $1
    AND cm.deleted_at IS NULL
WHERE c.deleted_at IS NULL
AND (
    c.group_id IS NULL
   
)
GROUP BY c.id, c.group_id, c.created_at, c.updated_at
ORDER BY c.updated_at DESC
LIMIT $2 OFFSET $3;




-- name: DeleteConversationByExactMembers :one
-- Delete a conversation only if its members exactly match the provided list.
-- Returns 0 rows if conversation doesn't exist, members don’t match exactly, conversation has extra or missing members.
WITH target_members AS (
    SELECT unnest(@member_ids::bigint[]) AS user_id
),
matched_conversation AS (
    SELECT cm.conversation_id
    FROM conversation_members cm
    JOIN target_members tm ON tm.user_id = cm.user_id
    GROUP BY cm.conversation_id
    HAVING 
        -- same count of overlapping members
        COUNT(*) = (SELECT COUNT(*) FROM target_members)
        -- and the conversation has no extra members
        AND COUNT(*) = (
            SELECT COUNT(*) 
            FROM conversation_members cm2 
            WHERE cm2.conversation_id = cm.conversation_id
              AND cm2.deleted_at IS NULL
        )
)
UPDATE conversations c
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE c.id = (SELECT conversation_id FROM matched_conversation)
RETURNING *;


