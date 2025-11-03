package models

import (
	"context"
	"time"
)

// /comments?parent_id=x&last_id=x&range=x
type CommentsRequest struct { //GET
	ParentID int64
	LastId   int64
	Range    int
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
	HaveMore bool      `json:"have_more_comments"`
}

type CommentsDbRequest struct {
	Ctx context.Context
	CommentsRequest
}

type CommentsDbResponse struct {
	CommentsResponse
}

// Create new Comment
// /new_comment {method: 'POST', body: {parent_id: int64, body: string}}
type NewCommentRequest struct {
	ParentId int64  `json:"parent_id"`
	Body     string `json:"body"`
}

type NewCommentResponse struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type NewCommentDbRequest struct {
	Ctx    context.Context
	UserId int64
	NewCommentRequest
}

type NewCommentDbResponse struct {
	NewCommentResponse
}

type CommentCreator struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"` // online, offline
}

type Comment struct {
	Id                   int64          `json:"id"`
	ParentId             int64          `json:"parent_id"`
	Body                 string         `json:"body"`
	ReactionCount        map[string]int `json:"reaction_count"`
	Creator              CommentCreator `json:"creator"`
	CurrentUserReactions []string       `json:"current_user_reaction"`
	CreatedAt            *time.Time     `json:"created_at"`
}
