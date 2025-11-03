package database

import "time"

type Users struct {
	Id        int64
	Username  string
	Gender    string
	FirstName string
	LastName  string
	Email     string
	Age       string
	Avatar    string
}

type Posts struct {
	Id           int64
	PostTitle    string
	PostBody     string
	PostReaction *string
	PostCreator  int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Comments struct {
	Id              int
	PostId          int
	CommentBody     string
	CommentReaction *string
	CommentCreator  int
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
