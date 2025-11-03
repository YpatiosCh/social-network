package models

import (
	"context"
	"time"
)

// /dm?last_id=x&before=10&after=10 // including last_id
type MessagesRequest struct {
	ConversationId int64
	LastId         int64 // return this
	Before         int64 // how many before LastId
	After          int64 // how many after LastId
}

type MessagesResponse struct {
	Messages       []Message `json:"messages"`
	HaveMoreBefore bool      `json:"have_more_dms_before"`
	HaveMoreAfter  bool      `json:"have_more_dms_after"`
}

type MessagesDbRequest struct {
	UserId int64
	Ctx    context.Context
	MessagesRequest
}

type MessagesDbResponse struct {
	MessagesResponse
}

// Create New message
/*
/new_dm {
    method: 'POST',
    body: { conversation_id: int64, receiver_id: int64, body: string }
    }
*/

// if convId then create in the handler the conversation
type NewMessageRequest struct {
	ReceiverId     int64  `json:"receiver_id"`
	ConversationId int64  `json:"conversation_id"`
	Body           string `json:"body"`
}

type NewMessageResponse struct {
	MessageId      int64     `json:"message_id"`
	CreatedAt      time.Time `json:"created_at"`
	ConversationId int64     `json:"conversation_id"`
}

type NewMessageDbRequest struct {
	Ctx      context.Context
	SenderId int64
	NewMessageRequest
}

type NewMessageDbResponse struct {
	NewMessageResponse
}

type Message struct {
	Id                   int64          `json:"id"`
	ConversationId       int64          `json:"conversation_id"`
	SenderId             int64          `json:"sender_id"`
	Body                 string         `json:"body"`
	ReactionCount        map[string]int `json:"reaction_count"`
	CurrentUserReactions []string       `json:"current_user_reactions"`
	CreatedAt            *time.Time     `json:"created_at"`
	Delivered            *bool          `json:"delivered"`
}

type MessageSender struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"` // online, offline
}

type TypingRequest struct {
	RecipientId int64 `json:"recipient_id"`
}
