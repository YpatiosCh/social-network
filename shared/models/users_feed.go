package models

import "context"

// /users_feed?last_id=x&range=10&with_conv=x
type UsersFeedRequest struct { //GET
	LastId int64
	Range  int
	/*
		If true I only want users that have a converation
		with me sorted by last interaction (recent first).
		Else bring users alphabetically
	*/
}

type UsersFeedDbRequest struct {
	Ctx    context.Context
	UserId int64
	UsersFeedRequest
}

type UsersFeedResponse struct {
	Feed     []Feed `json:"users_feed"`
	HaveMore bool   `json:"have_more_users"`
}

type UsersFeedDbResponse struct {
	UsersFeedResponse
}

type Feed struct {
	User                UserOnFeed          `json:"user"`                           // id,username, avatar, status
	ConversationDetails ConversationDetails `json:"conversation_details,omitempty"` // null or ConversationId, UnreadCount , LastReadMessageId
}

type ConversationDetails struct {
	ConversationId    uint64  `json:"id"`
	UnreadCount       int     `json:"unread_count"`         // the logged in users unread count (exluding own messages)
	LastReadMessageId *uint64 `json:"last_read_message_id"` // the logged in user's last read (even if the creator is the logged in user)
}

type UserOnFeed struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"` // online, offline
}
