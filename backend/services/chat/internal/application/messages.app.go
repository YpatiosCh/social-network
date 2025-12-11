package application

import (
	"context"
	"fmt"
	"social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation_id and user_id fails returns error.
func (c *ChatService) CreateMessage(ctx context.Context,
	params md.CreateMessageParams) (msg md.MessageResp, err error) {
	if err := customtypes.ValidateStruct(params); err != nil {
		return msg, err
	}
	if (msg == md.MessageResp{}) {
		return msg, fmt.Errorf("user is not a member of conversation id: %v", params.ConversationId)
	}
	return c.Queries.CreateMessage(ctx, params)
}

func (c *ChatService) GetMessages(ctx context.Context,
	params md.GetMessagesParams) (messages []md.MessageResp, err error) {
	if err := customtypes.ValidateStruct(params); err != nil {
		return messages, err
	}
	return c.Queries.GetMessages(ctx, params)
}

func (c *ChatService) UpdateLastReadMessage(ctx context.Context,
	params md.UpdateLastReadMessageParams,
) (member md.ConversationMember, err error) {
	if err := customtypes.ValidateStruct(params); err != nil {
		return member, err
	}
	if (member == md.ConversationMember{}) {
		return member, fmt.Errorf("user is not a member of conversation id: %v", params.ConversationId)
	}
	return c.Queries.UpdateLastReadMessage(ctx, params)
}
