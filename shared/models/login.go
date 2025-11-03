package models

// /login
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	User LoginResponseUser `json:"user"`
}

type LoginResponseUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type LoginDbRequest struct {
	LoginRequest
}

type LoginDbResponse struct {
	LoginResponseUser
}
