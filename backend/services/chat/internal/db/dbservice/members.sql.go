package dbservice

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

const getConversationMembers = `-- name: GetConversationMembers :many
SELECT cm2.user_id
FROM conversation_members cm1
JOIN conversation_members cm2
  ON cm2.conversation_id = cm1.conversation_id
WHERE cm1.user_id = $2
  AND cm2.conversation_id = $1
  AND cm2.user_id <> $2
  AND cm2.deleted_at IS NULL
`

// Returns memebers of a conversation that user is a member.
// OK!
func (q *Queries) GetConversationMembers(ctx context.Context,
	arg md.GetConversationMembersParams) (members ct.Ids, err error) {
	rows, err := q.db.Query(ctx,
		getConversationMembers,
		arg.ConversationId,
		arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members = ct.Ids{}
	for rows.Next() {
		var user_id int64
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		members = append(members, ct.Id(user_id))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

const deleteConversationMember = `
UPDATE conversation_members to_delete
SET deleted_at = NOW()
FROM conversation_members owner
WHERE to_delete.conversation_id = $1
  AND to_delete.user_id = $2
  AND to_delete.deleted_at IS NULL
  AND owner.conversation_id = $1
  AND owner.user_id = $3
  AND owner.deleted_at IS NULL
RETURNING to_delete.conversation_id, to_delete.user_id, to_delete.last_read_message_id, cm_target.joined_at, cm_target.deleted_at
`

// Deletes conversation member from conversation where user tagged as owner is a part of.
// Returnes user deleted details. If no rows returned means no deletation occured.
// Can be used for self deletation if owner and toDelete are that same id.
func (q *Queries) DeleteConversationMember(ctx context.Context,
	arg md.DeleteConversationMemberParams,
) (dltMember md.ConversationMemberDeleted, err error) {
	row := q.db.QueryRow(ctx, deleteConversationMember,
		arg.ConversationID, arg.ToDelete, arg.Owner)
	err = row.Scan(
		&dltMember.ConversationId,
		&dltMember.UserId,
		&dltMember.LastReadMessageId,
		&dltMember.JoinedAt,
		&dltMember.DeletedAt,
	)
	return dltMember, err
}
