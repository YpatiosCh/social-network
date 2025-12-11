package dbservice

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ConversationMember struct {
	ConversationID    int64
	UserID            int64
	LastReadMessageID pgtype.Int8
	JoinedAt          pgtype.Timestamptz
	DeletedAt         pgtype.Timestamptz
}

type Message struct {
	ID             int64
	ConversationID int64
	SenderID       pgtype.Int8
	MessageText    string
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
	DeletedAt      pgtype.Timestamptz
}
