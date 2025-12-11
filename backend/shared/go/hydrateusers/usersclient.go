package userhydrate

import (
	"context"
	userpb "social-network/shared/gen-go/users"
)

// UsersBatchClient is the subset the hydrator needs.
type UsersBatchClient interface {
	GetBatchBasicUserInfo(ctx context.Context, userIds []int64) (*userpb.ListUsers, error)
}
