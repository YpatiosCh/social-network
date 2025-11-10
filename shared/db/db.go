package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectOrCreateDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connStr := cfg.ConnString()
	log.Printf("Connecting to database: %s", connStr)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New failed: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return pool, nil
}
