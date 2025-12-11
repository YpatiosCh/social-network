package models

type HasUser interface {
	GetUserId() int64
	SetUser(User)
}
