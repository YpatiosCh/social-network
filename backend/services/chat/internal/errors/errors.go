package errors

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInternal      = errors.New("internal error")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalid       = errors.New("invalid arguments")
	ErrUsersService  = errors.New("users service error")
)
