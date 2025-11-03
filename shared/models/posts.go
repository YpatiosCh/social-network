package models

import (
	"context"
	"time"
)

// /posts?last_id=x&range=10&categories=x
type PostsFeedRequest struct { // GET
	LastSeenId int64
	Range      int
	Categories []int64 // comma seperated list of ids "2,3"
}

type PostsFeedResponse struct {
	Posts    []Post `json:"posts"`
	HaveMore bool   `json:"have_more_posts"`
}

type PostsFeedDbRequest struct {
	Ctx    context.Context
	UserId int64
	PostsFeedRequest
}

type PostsFeedDbResponse struct {
	PostsFeedResponse
}

// Create New Post
// /new_post  {method: 'POST', body: { title: string, body: string } }
type NewPostRequest struct {
	Title      string   `json:"title"`
	Body       string   `json:"body"`
	Categories []string `json:"categories"` // comma seperated list of ids "2,3"
}

type NewPostResponse struct {
	PostId    int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type NewPostDbRequest struct {
	Ctx    context.Context
	UserId int64
	NewPostRequest
}

type NewPostDbResponse struct {
	NewPostResponse
}

// /post?id=x
type SinglePostRequest struct { // GET
	Id int64
}

type SinglePostResponse struct {
	*Post
}

type SinglePostDbRequest struct { // GET
	SinglePostRequest
}

type SinglePostDbResponse struct {
	*Post
}

type Post struct {
	Id                   int64          `json:"id"`
	Creator              PostCreator    `json:"creator"`
	CreatedAt            *time.Time     `json:"created_at"`
	Title                string         `json:"title"`
	Body                 string         `json:"body"`
	ReactionCount        map[string]int `json:"reaction_count"` // key = reaction emoji
	CurrentUserReactions []string       `json:"current_user_reactions"`
	CommentCount         int            `json:"comment_count"`
	Categories           []Category     `json:"categories"`
	Truncated            bool           `json:"truncated"`
	FirstCommentId       int64          `json:"first_comment_id"` // only on full post
}

type PostCreator struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
