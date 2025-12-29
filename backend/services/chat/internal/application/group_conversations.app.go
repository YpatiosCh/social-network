package application

import (
	"context"
	"fmt"
	ce "social-network/services/chat/internal/errors"
	ct "social-network/shared/go/ct"
	md "social-network/shared/go/models"
)

func (c *ChatService) AddMembersToGroupConversation(ctx context.Context,
	arg md.AddMembersToGroupConversationParams) (convId ct.Id, err error) {
	errMsg := fmt.Sprintf("add memvbers to group conversation: arg: %#v", arg)

	if err := ct.ValidateStruct(arg); err != nil {
		return convId, ct.Wrap(ce.ErrInvalid, err, errMsg)
	}

	convId, err = c.Queries.AddMembersToGroupConversation(ctx, arg)
	if err != nil {
		return convId, ct.Wrap(nil, err, errMsg)
	}
	return convId, nil
}
