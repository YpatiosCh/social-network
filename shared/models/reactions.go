package models

import "context"

/*
	/new_reaction {
	     method: "POST",
	     body: { content_id: int64, type: string, new: bool }
	}
*/
type UserReacted struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type NewReactionRequest struct {
	ContentId    int64  `json:"content_id"`
	ReactionType string `json:"type"`
	New          bool   `json:"new"`
}

type NewReactionDbRequest struct {
	Ctx    context.Context
	UserId int64
	NewReactionRequest
}

// OPTIONAL
// /usersreacted?content_id=x&type=x&last_id=x&range=10
type UsersReactedRequest struct { // Get
	ContentId    int64  // comment or post
	ReactionType string // emoji
	LastID       int64  // the last user that I see
	Range        int
}

type UsersReactedResponse struct {
	Users []UserReacted `json:"users_reacted"`
	Total int64         `json:"total"`
}

type UsersReactedDbRequest struct {
	Ctx context.Context
	UsersReactedRequest
}

type UsersReactedDbResponse struct {
	UsersReactedResponse
}

type DeleteReactionRequest struct {
	PostID int64 `json:"post_id"`
}
