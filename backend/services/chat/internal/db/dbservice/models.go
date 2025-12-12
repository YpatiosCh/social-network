package dbservice

import (
	ct "social-network/shared/go/customtypes"
)

// Conversations
type GetUserConversationsRow struct {
	ConversationId       ct.Id
	CreatedAt            ct.GenDateTime
	UpdatedAt            ct.GenDateTime
	MemberIds            ct.Ids
	UnreadCount          int64
	FirstUnreadMessageId ct.Id `validation:"nullable"`
}
