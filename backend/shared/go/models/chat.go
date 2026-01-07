package models

import (
	ct "social-network/shared/go/ct"
)

type GetOrCreatePCRec struct {
	User              ct.Id `json:"user"`
	OtherUser         ct.Id `json:"other_user"`
	RetrieveOtherUser bool  `json:"retrieve_other_user"`
}

type GetOrCreatePCResp struct {
	ConversationId  ct.Id
	OtherUser       User
	LastReadMessage ct.Id `validation:"nullable"`
	IsNew           bool
}

type CreateGCParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"user_ids"`
}

type CreateGMParams struct {
	GroupId     ct.Id      `json:"group_id"`
	SenderId    ct.Id      `json:"sender_id"`
	MessageText ct.MsgBody `json:"message_text"`
}

type CreatePMParams struct {
	ConversationId ct.Id      `json:"conversation_id"`
	SenderId       ct.Id      `json:"sender_id"`
	MessageText    ct.MsgBody `json:"message_text"`
}

type GetPMsParams struct {
	ConversationId    ct.Id    `json:"conversation_id"`
	UserId            ct.Id    `json:"user_id"`
	BoundaryMessageId ct.Id    `json:"boundary_message_id" validation:"nullable"`
	Limit             ct.Limit `json:"limit"`
	RetrieveUsers     bool     `json:"retrieve_users"`
}

type GetPMsResp struct {
	HaveMore bool
	Messages []PM
}

type GetPCsReq struct {
	UserId     ct.Id          `json:"user_id"`
	BeforeDate ct.GenDateTime `json:"before_date"`
	Limit      ct.Limit       `json:"limit"`
}

type PCsPreview struct {
	ConversationId ct.Id
	UpdatedAt      ct.GenDateTime
	OtherUser      User
	LastMessage    PM
	UnreadCount    int
}

type PM struct {
	Id             ct.Id
	ConversationID ct.Id
	Sender         User
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime `validation:"nullable"`
	UpdatedAt      ct.GenDateTime `validation:"nullable"`
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type UpdateLastReadMsgParams struct {
	ConversationId    ct.Id `json:"conversation_id"`
	UserId            ct.Id `json:"user_id"`
	LastReadMessageId ct.Id `json:"last_read_message_id"`
}

type ConversationMember struct {
	ConversationId    ct.Id
	UserId            ct.Id
	LastReadMessageId ct.Id `validation:"nullable"`
	JoinedAt          ct.GenDateTime
	DeletedAt         ct.GenDateTime `validation:"nullable"`
}
