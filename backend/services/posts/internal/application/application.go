package application

import (
	"context"
	"social-network/services/posts/internal/client"
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	db       sqlc.Querier
	txRunner TxRunner
	clients  ClientsInterface
	hydrator *UserHydrator
}

type UserHydrator struct {
	clients *client.Clients
	//cache   RedisCache
}

// ClientsInterface defines the methods that Application needs from clients.
type ClientsInterface interface {
	IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error)
	IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error)
	GetFollowingIds(ctx context.Context, userId int64) ([]int64, error)
}

func NewUserHydrator(clients *client.Clients) *UserHydrator {
	return &UserHydrator{
		clients: clients,
	}
}

// NewApplication constructs a new Application with transaction support
func NewApplication(db sqlc.Querier, pool *pgxpool.Pool, clients *client.Clients) *Application {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := db.(*sqlc.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}

	return &Application{
		db:       db,
		txRunner: txRunner,
		clients:  clients,
		hydrator: NewUserHydrator(clients),
	}
}

func NewApplicationWithMocks(db sqlc.Querier, clients ClientsInterface) *Application {
	return &Application{
		db:      db,
		clients: clients,
	}
}
func NewApplicationWithMocksTx(db sqlc.Querier, clients ClientsInterface, txRunner TxRunner) *Application {
	return &Application{
		db:       db,
		clients:  clients,
		txRunner: txRunner,
		hydrator: NewUserHydrator(nil), // or pass clients if needed
	}
}
