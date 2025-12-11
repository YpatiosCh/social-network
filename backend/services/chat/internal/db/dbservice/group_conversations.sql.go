package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

const addMembersToGroupConversation = `-- name: AddMembersToGroupConversation :one
WITH convo AS (
    SELECT id
    FROM conversations
    WHERE group_id = $1
      AND deleted_at IS NULL
),
insert_members AS (
    INSERT INTO conversation_members (conversation_id, user_id)
    SELECT (SELECT id FROM convo), unnest($2::bigint[])
    ON CONFLICT (conversation_id, user_id) DO NOTHING
    RETURNING conversation_id
)
SELECT id FROM convo
`

// Find a conversation by group_id and insert the given user_ids into conversation_members.
// existing members are ignored, new members are added.
func (q *Queries) AddMembersToGroupConversation(ctx context.Context,
	arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	row := q.db.QueryRow(ctx,
		addMembersToGroupConversation,
		arg.GroupId,
		arg.UserIds,
	)
	err = row.Scan(&convId)
	return convId, err
}
