package application

import (
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsService struct {
	db   sqlc.Querier  // interface, can be *sqlc.Queries or mock
	pool *pgxpool.Pool // needed to start transactions
}

// NewPostsService constructs a new PostsService
func NewPostsService(db sqlc.Querier, pool *pgxpool.Pool) *PostsService {
	return &PostsService{
		db:   db,
		pool: pool,
	}
}
