package models

import (
	ct "social-network/shared/go/customtypes"
)

type CreatePrivateConvParams struct {
	UserA ct.Id `json:"user_a"`
	UserB ct.Id `json:"user_b"`
}

type CreateGroupConvParams struct {
	GroupId ct.Id  `json:"group_id"`
	UserIds ct.Ids `json:"users_id"`
}

type Conversation struct {
	ID        ct.Id
	GroupID   ct.Id
	CreatedAt ct.GenDateTime `validation:"nullable"`
	UpdatedAt ct.GenDateTime `validation:"nullable"`
	DeletedAt ct.GenDateTime `validation:"nullable"`
}

type AddMembersToGroupConversationParams struct {
	GroupID ct.Id
	UserIds ct.Ids
}
