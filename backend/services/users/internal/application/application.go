package application

import (
	"context"
	"social-network/services/users/internal/client"
	"social-network/services/users/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db      sqlc.Querier  // interface, can be *sqlc.Queries or mock
	pool    *pgxpool.Pool // needed to start transactions
	clients ClientsInterface
}

// NewApplication constructs a new UserService
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool, clients *client.Clients) *Application {
	return &Application{
		db:      db,
		pool:    pool,
		clients: clients,
	}
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	CreateGroupConversation(ctx context.Context, groupId int64, ownerId int64) error
	CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error
}
