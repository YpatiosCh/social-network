package application

import (
	"social-network/services/chat/internal/client"
	"social-network/services/chat/internal/db/dbservice"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Holds logic for requests and calls
type ChatService struct {
	Clients  *client.Clients
	Queries  dbservice.Querier
	txRunner TxRunner
}

func NewChatService(pool *pgxpool.Pool, clients *client.Clients, queries dbservice.Querier) *ChatService {
	var txRunner TxRunner
	if pool != nil {
		queries, ok := queries.(*dbservice.Queries)
		if !ok {
			panic("db must be *sqlc.Queries for transaction support")
		}
		txRunner = NewPgxTxRunner(pool, queries)
	}
	return &ChatService{
		Clients:  clients,
		Queries:  queries,
		txRunner: txRunner,
	}
}
