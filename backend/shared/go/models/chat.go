package models

import (
	ct "social-network/shared/go/customtypes"
)

type AddConversationMembersParams struct {
	ConversationId ct.Id
	UserIds        ct.Ids
}

type AddMembersToGroupConversationParams struct {
	GroupId ct.Id
	UserIds ct.Ids
}

type CreatePrivateConvParams struct {
	UserA ct.Id `json:"user_a"`
	UserB ct.Id `json:"user_b"`
}

type CreateGroupConvParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"users_id"`
}

type CreateMessageParams struct {
	ConversationId ct.Id
	SenderId       ct.Id
	MessageText    ct.MsgBody
}

type GetNextMessageParams struct {
	FirstMessageId ct.Id
	ConversationId ct.Id
	UserId         ct.Id
	Limit          ct.Limit
	Offset         ct.Offset
}

type ConversationDeleteResp struct {
	Id        ct.Id
	GroupId   ct.Id
	CreatedAt ct.GenDateTime
	UpdatedAt ct.GenDateTime
	DeletedAt ct.GenDateTime
}

type ConversationResponse struct {
	Id             ct.Id
	GroupId        ct.Id
	LastMessageId  ct.Id
	FirstMessageId ct.Id
	CreatedAt      ct.GenDateTime
	UpdatedAt      ct.GenDateTime `validation:"nullable"`
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type ConversationMember struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime `validation:"nullable"`
}

// All fields are required except LastReadMessgeId
type ConversationMemberDeleted struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime
}

type GetConversationMembersParams struct {
	ConversationId ct.Id
	UserID         ct.Id
}

type GetPrevMessagesParams struct {
	UserId            ct.Id
	ConversationId    ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	Limit             ct.Limit
	Offset            ct.Offset
}

type GetPrevMessagesResp struct {
	FirstMessageId ct.Id
	HaveMore       bool
	Messages       []MessageResp
}

type GetUserConversationsParams struct {
	UserId  ct.Id
	GroupId ct.Id
	Limit   ct.Limit
	Offset  ct.Offset
}

type GetUserConversationsRow struct {
	ConversationId       ct.Id
	CreatedAt            ct.GenDateTime
	UpdatedAt            ct.GenDateTime
	MemberIds            []int64
	UnreadCount          int64
	FirstUnreadMessageId *int64
}

// All fields are required except deleted at which in most cases is null.
type MessageResp struct {
	Id             ct.Id
	ConversationID ct.Id
	SenderID       ct.Id
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime
	UpdatedAt      ct.GenDateTime
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type DeleteConversationMemberParams struct {
	ConversationID ct.Id
	Owner          ct.Id
	ToDelete       ct.Id
}

// Last Read message is not nullable. If it is null then request is invalid.
type UpdateLastReadMessageParams struct {
	ConversationId    ct.Id
	UserID            ct.Id
	LastReadMessageId ct.Id
}
