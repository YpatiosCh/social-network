package dbservice

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid request")
	ErrInternal = errors.New("internal")
)
