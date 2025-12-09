package mocks

import (
	"context"
	"social-network/services/posts/internal/db/sqlc"

	"github.com/stretchr/testify/mock"
)

// MockTxRunner is a simple mock usable in tests. It can be given a Querier
// instance which will be passed to the transactional function when RunTx is
// called. Callers can set the Querier field to a *mocks.MockQuerier instance.
type MockTxRunner struct {
	mock.Mock
	Querier sqlc.Querier
}

func (m *MockTxRunner) RunTx(ctx context.Context, fn func(sqlc.Querier) error) error {
	m.Called(ctx)
	if m.Querier != nil {
		return fn(m.Querier)
	}
	return nil
}
