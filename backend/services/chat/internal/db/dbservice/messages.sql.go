package dbservice

import (
	"context"
	md "social-network/shared/go/models"
)

// Updates the given users last read message in given conversation to given message id.
// FIX
func (q *Queries) UpdateLastReadMessage(ctx context.Context,
	arg md.UpdateLastReadMsgParams,
) (member md.ConversationMember, err error) {
	row := q.db.QueryRow(ctx, updateLastReadMessage, arg.ConversationId, arg.UserId, arg.LastReadMessageId)
	err = row.Scan(
		&member.ConversationId,
		&member.UserId,
		&member.LastReadMessageId,
		&member.JoinedAt,
		&member.DeletedAt,
	)
	return member, err
}
