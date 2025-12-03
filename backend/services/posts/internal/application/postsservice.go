package application

import (
	"social-network/services/posts/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsService struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

// NewPostsService constructs a new PostsService
func NewPostsService(pool *pgxpool.Pool) *PostsService {
	return &PostsService{
		q:    sqlc.New(pool),
		pool: pool,
	}
}
