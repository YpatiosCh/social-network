package application_test

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"social-network/services/posts/internal/application"
// 	clientmocks "social-network/services/posts/internal/client/mocks"
// 	mocks "social-network/services/posts/internal/db/mocks"
// 	"social-network/services/posts/internal/db/sqlc"
// 	ct "social-network/shared/go/customtypes"

// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // ===============================
// // TEST SETUP HELPERS
// // ===============================

// func setupTestApp() (*application.Application, *mocks.MockQuerier, *clientmocks.MockClients) {
// 	mockDB := &mocks.MockQuerier{}
// 	mockClients := &clientmocks.MockClients{}
// 	txRunner := &mockTxRunner{db: mockDB}

// 	app := application.NewApplicationWithMocksTx(mockDB, mockClients, txRunner)
// 	return app, mockDB, mockClients
// }

// // mockTxRunner implements TxRunner for testing
// type mockTxRunner struct {
// 	db sqlc.Querier
// }

// func (m *mockTxRunner) RunTx(ctx context.Context, fn func(q sqlc.Querier) error) error {
// 	return fn(m.db)
// }

// // ===============================
// // CREATE POST
// // ===============================

// func TestCreatePost(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		req           application.CreatePostReq
// 		setupMock     func(*mocks.MockQuerier, *clientmocks.MockClients)
// 		expectedError error
// 	}{
// 		{
// 			name: "successful public post",
// 			req: application.CreatePostReq{
// 				Body:      ct.PostBody("Hello World"),
// 				CreatorId: ct.Id(1),
// 				GroupId:   ct.Id(10),
// 				Audience:  ct.Audience("everyone"),
// 				Image:     ct.Id(5),
// 			},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				c.On("IsGroupMember", mock.Anything, int64(1), int64(10)).Return(true, nil)
// 				m.On("CreatePost", mock.Anything, mock.MatchedBy(func(p sqlc.CreatePostParams) bool {
// 					return p.PostBody == "Hello World" && p.CreatorID == 1
// 				})).Return(int64(100), nil)
// 				m.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "user not member of group",
// 			req: application.CreatePostReq{
// 				Body:      ct.PostBody("Hello World"),
// 				CreatorId: ct.Id(1),
// 				GroupId:   ct.Id(10),
// 				Audience:  ct.Audience("group"),
// 			},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				c.On("IsGroupMember", mock.Anything, int64(1), int64(10)).Return(false, nil)
// 			},
// 			expectedError: application.ErrNotAllowed,
// 		},
// 		{
// 			name: "selected audience with IDs",
// 			req: application.CreatePostReq{
// 				Body:        ct.PostBody("Hello"),
// 				CreatorId:   ct.Id(2),
// 				GroupId:     ct.Id(0),
// 				Audience:    ct.Audience("selected"),
// 				AudienceIds: ct.Ids{1, 2},
// 			},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				m.On("CreatePost", mock.Anything, mock.Anything).Return(int64(101), nil)
// 				m.On("InsertPostAudience", mock.Anything, mock.Anything).Return(int64(2), nil)
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			app, mockDB, mockClients := setupTestApp()
// 			tt.setupMock(mockDB, mockClients)

// 			err := app.CreatePost(context.Background(), tt.req)

// 			if tt.expectedError != nil {
// 				assert.ErrorIs(t, err, tt.expectedError)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockDB.AssertExpectations(t)
// 			mockClients.AssertExpectations(t)
// 		})
// 	}
// }

// // ===============================
// // DELETE POST
// // ===============================

// func TestDeletePost(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		req           application.GenericReq
// 		setupMock     func(*mocks.MockQuerier, *clientmocks.MockClients)
// 		expectedError error
// 	}{
// 		{
// 			name: "successful deletion",
// 			req:  application.GenericReq{EntityId: ct.Id(100), RequesterId: ct.Id(1)},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				m.On("GetEntityCreatorAndGroup", mock.Anything, int64(100)).
// 					Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 1, GroupID: 10}, nil)
// 				c.On("IsFollowing", mock.Anything, int64(1), int64(1)).Return(false, nil)
// 				c.On("IsGroupMember", mock.Anything, int64(1), int64(10)).Return(true, nil)
// 				m.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 				m.On("DeletePost", mock.Anything, mock.Anything).Return(int64(1), nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "post not found",
// 			req:  application.GenericReq{EntityId: ct.Id(999), RequesterId: ct.Id(1)},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				m.On("GetEntityCreatorAndGroup", mock.Anything, int64(999)).
// 					Return(sqlc.GetEntityCreatorAndGroupRow{CreatorID: 1, GroupID: 10}, nil)
// 				c.On("IsFollowing", mock.Anything, int64(1), int64(1)).Return(false, nil)
// 				c.On("IsGroupMember", mock.Anything, int64(1), int64(10)).Return(true, nil)
// 				m.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 				m.On("DeletePost", mock.Anything, mock.Anything).Return(int64(0), nil)
// 			},
// 			expectedError: application.ErrNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			app, mockDB, mockClients := setupTestApp()
// 			tt.setupMock(mockDB, mockClients)

// 			err := app.DeletePost(context.Background(), tt.req)
// 			if tt.expectedError != nil {
// 				assert.ErrorIs(t, err, tt.expectedError)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockDB.AssertExpectations(t)
// 			mockClients.AssertExpectations(t)
// 		})
// 	}
// }

// // ===============================
// // EDIT POST
// // ===============================

// func TestEditPost(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		req           application.EditPostReq
// 		setupMock     func(*mocks.MockQuerier, *clientmocks.MockClients)
// 		expectedError error
// 	}{
// 		{
// 			name: "edit body and image",
// 			req: application.EditPostReq{
// 				PostId:      ct.Id(100),
// 				RequesterId: ct.Id(1),
// 				NewBody:     ct.PostBody("Updated"),
// 				Image:       ct.Id(5),
// 				Audience:    "everyone",
// 			},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				m.On("EditPostContent", mock.Anything, mock.Anything).Return(int64(1), nil)
// 				m.On("UpsertImage", mock.Anything, mock.Anything).Return(nil)
// 				m.On("UpdatePostAudience", mock.Anything, mock.Anything).Return(int64(1), nil)
// 				c.On("IsFollowing", mock.Anything, int64(1), int64(1)).Return(false, nil)
// 				c.On("IsGroupMember", mock.Anything, int64(1), int64(10)).Return(true, nil)
// 				m.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "edit selected audience with ids",
// 			req: application.EditPostReq{
// 				PostId:      ct.Id(101),
// 				RequesterId: ct.Id(2),
// 				NewBody:     ct.PostBody("New"),
// 				Audience:    "selected",
// 				AudienceIds: ct.Ids{3, 4},
// 			},
// 			setupMock: func(m *mocks.MockQuerier, c *clientmocks.MockClients) {
// 				m.On("ClearPostAudience", mock.Anything, int64(101)).Return(nil)
// 				m.On("InsertPostAudience", mock.Anything, mock.Anything).Return(int64(2), nil)
// 				m.On("UpdatePostAudience", mock.Anything, mock.Anything).Return(int64(1), nil)
// 				m.On("EditPostContent", mock.Anything, mock.Anything).Return(int64(1), nil)
// 				c.On("IsFollowing", mock.Anything, int64(2), int64(2)).Return(false, nil)
// 				c.On("IsGroupMember", mock.Anything, int64(2), int64(10)).Return(true, nil)
// 				m.On("CanUserSeeEntity", mock.Anything, mock.Anything).Return(true, nil)
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			app, mockDB, mockClients := setupTestApp()
// 			tt.setupMock(mockDB, mockClients)

// 			err := app.EditPost(context.Background(), tt.req)
// 			if tt.expectedError != nil {
// 				assert.ErrorIs(t, err, tt.expectedError)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockDB.AssertExpectations(t)
// 			mockClients.AssertExpectations(t)
// 		})
// 	}
// }

// // ===============================
// // GET MOST POPULAR POST
// // ===============================

// func TestGetMostPopularPostInGroup(t *testing.T) {
// 	now := time.Now()
// 	tests := []struct {
// 		name      string
// 		groupID   int64
// 		setupMock func(*mocks.MockQuerier)
// 		expectErr error
// 	}{
// 		{
// 			name:    "post found",
// 			groupID: 10,
// 			setupMock: func(m *mocks.MockQuerier) {
// 				m.On("GetMostPopularPostInGroup", mock.Anything, mock.Anything).
// 					Return(sqlc.GetMostPopularPostInGroupRow{
// 						ID:              1,
// 						PostBody:        "Hi",
// 						CreatorID:       1,
// 						Audience:        "everyone",
// 						CommentsCount:   3,
// 						ReactionsCount:  5,
// 						LastCommentedAt: pgtype.Timestamptz{Time: now, Valid: true},
// 						CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
// 						UpdatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
// 						Image:           5,
// 					}, nil)
// 			},
// 			expectErr: nil,
// 		},
// 		{
// 			name:    "no posts",
// 			groupID: 99,
// 			setupMock: func(m *mocks.MockQuerier) {
// 				m.On("GetMostPopularPostInGroup", mock.Anything, mock.Anything).
// 					Return(sqlc.GetMostPopularPostInGroupRow{}, errors.New("sql: no rows in result set"))
// 			},
// 			expectErr: application.ErrNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			app, mockDB, _ := setupTestApp()
// 			tt.setupMock(mockDB)

// 			post, err := app.GetMostPopularPostInGroup(context.Background(), application.SimpleIdReq{Id: ct.Id(tt.groupID)})

// 			if tt.expectErr != nil {
// 				assert.ErrorIs(t, err, tt.expectErr)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, ct.Id(1), post.PostId)
// 				assert.Equal(t, ct.PostBody("Hi"), post.Body)
// 				assert.Equal(t, ct.Id(5), post.Image)
// 			}

// 			mockDB.AssertExpectations(t)
// 		})
// 	}
// }
