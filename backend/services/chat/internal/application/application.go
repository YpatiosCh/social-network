package application

import (
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type ChatService struct {
	Pool    *pgxpool.Pool
	Clients *client.Clients
	Queries sqlc.Querier
}
