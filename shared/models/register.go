package models

import "context"

// /register {POST}
type RegisterRequest struct {
	Username        *string `json:"username"` // required
	Email           *string `json:"email"`    // optional if you allow email login
	Password        string  `json:"password"` // required
	ConfirmPassword string  `json:"confirm_password"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Gender          *string `json:"gender"`
	Age             *string `json:"age"` // prefer int or a Birthdate
	Avatar          string  `json:"avatar"`
}

type RegisterResponse struct {
	User RegisterResponseUser `json:"user"`
}

type RegisterResponseUser struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}

type RegisterDbRequest struct {
	Ctx          context.Context
	Username     *string
	Email        *string
	FirstName    string
	LastName     string
	Gender       *string
	Age          *string
	Avatar       string
	Identifier   string
	PasswordHash string
}

type RegisterDbResponse struct {
	RegisterResponseUser
}
