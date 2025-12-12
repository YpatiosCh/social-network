package application

import (
	"context"
	"fmt"
	"social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Creates a message row with conversation id if user is a memeber.
// If user match of conversation_id and user_id fails returns error.
//
// TODO: Hydrate and return full response
func (c *ChatService) CreateMessage(ctx context.Context,
	params md.CreateMessageParams) (msg md.MessageResp, err error) {
	if err := customtypes.ValidateStruct(params); err != nil {
		return msg, err
	}
	if (msg == md.MessageResp{}) {
		return msg, fmt.Errorf("user is not a member of conversation id: %v", params.ConversationId)
	}

	_, err = c.Queries.CreateMessage(ctx, params)
	if err != nil {
		return msg, err
	}
	return msg, err
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
