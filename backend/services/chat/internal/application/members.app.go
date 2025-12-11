package application

import (
	"context"
	ct "social-network/shared/go/customtypes"
	md "social-network/shared/go/models"
)

// Returns memebers of a conversation that user is a member.
func (c *ChatService) GetConversationMembers(ctx context.Context,
	params md.GetConversationMembersParams) (members ct.Ids, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return members, err
	}
	return c.Queries.GetConversationMembers(ctx, params)
}

func (c *ChatService) DeleteConversationMember(ctx context.Context,
	params md.DeleteConversationMemberParams,
) (dltMember md.ConversationMemberDeleted, err error) {
	if err := ct.ValidateStruct(params); err != nil {
		return dltMember, err
	}
	return c.Queries.DeleteConversationMember(ctx, params)
}
