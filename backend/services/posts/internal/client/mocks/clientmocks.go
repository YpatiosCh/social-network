package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClients struct {
	mock.Mock
}

func (m *MockClients) IsFollowing(ctx context.Context, userId, targetUserId int64) (bool, error) {
	ret := m.Called(ctx, userId, targetUserId)
	return ret.Bool(0), ret.Error(1)
}

func (m *MockClients) IsGroupMember(ctx context.Context, userId, groupId int64) (bool, error) {
	ret := m.Called(ctx, userId, groupId)
	return ret.Bool(0), ret.Error(1)
}

func (m *MockClients) GetFollowingIds(ctx context.Context, userId int64) ([]int64, error) {
	ret := m.Called(ctx, userId)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).([]int64), ret.Error(1)
}
