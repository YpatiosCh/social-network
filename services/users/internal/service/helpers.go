package userservice

import (
	"context"
	"fmt"
	"social-network/services/users/internal/db/sqlc"
)

// runTx runs a function inside a database transaction.
// If fn returns an error, the tx is rolled back.
func (s *UserService) runTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	// start tx
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// db must be a *sqlc.Queries to use WithTx
	base, ok := s.db.(*sqlc.Queries)
	if !ok {
		return fmt.Errorf("UserService.db must be *sqlc.Queries for transactions")
	}

	qtx := base.WithTx(tx)

	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
