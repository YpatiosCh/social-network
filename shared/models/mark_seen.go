package models

import "context"

// /mark-seen, { method: 'POST',  body: "dm_id": int}
type MarkSeenRequest struct {
	DmId           int64 `json:"dm_id"` // this is the last seen dm from active user on conversation
	ConversationId int64 `json:"conversation_id"`
}

type MarkSeenDbRequest struct {
	Ctx    context.Context
	UserId int64
	MarkSeenRequest
}
