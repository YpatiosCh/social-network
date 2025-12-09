package client

import (
	"context"
	chatpb "social-network/shared/gen-go/chat"
)

// Holds connections to clients
type Clients struct {
	ChatClient chatpb.ChatServiceClient
}

// on successful follow (public profile or accept follow request)
func (c *Clients) CreatePrivateConversation(ctx context.Context, userId1, userId2 int64) error {
	_, err := c.ChatClient.CreatePrivateConversation(ctx, &chatpb.CreatePrivateConvParams{
		UserA: userId1,
		UserB: userId2,
	})
	if err != nil {
		return err
	}
	return nil
}

// when group is created there's only the owner
func (c *Clients) CreateGroupConversation(ctx context.Context, groupId, ownerId int64) error {
	_, err := c.ChatClient.CreateGroupConversation(ctx, &chatpb.CreateGroupConvParams{
		GroupId: groupId,
		UserIds: []int64{ownerId},
	})
	if err != nil {
		return err
	}
	return nil
}

//Function to deactivate PrivateConversation on unfollow if !isfollowingeither

//Function to add group members?
