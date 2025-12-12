// Internal Models contains raw db responses
// before these are compiled by the application layer
package internalmodels

import ct "social-network/shared/go/customtypes"

// MESSAGES

type Message struct {
	Id             ct.Id
	ConversationID ct.Id
	SenderID       ct.Id
	MessageText    ct.MsgBody
	CreatedAt      ct.GenDateTime
	UpdatedAt      ct.GenDateTime
	DeletedAt      ct.GenDateTime `validation:"nullable"`
}

type GetPrevMessagesResp struct {
	FirstMessageId ct.Id
	Messages       []Message
}

type GetNextMessagesResp struct {
	FirstMessageId ct.Id
	Messages       []Message
}

// MEMBERS
