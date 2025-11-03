package models

import "context"

// Response to /user. User Id is taken from token
type UserResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type UserDbRequest struct {
	Ctx    context.Context
	UserId int64
}

type UserDbResponse struct {
	UserResponse
}
