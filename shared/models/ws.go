package models

type WSMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type ContentStatusUpdateWS struct {
	Id                   int64          `json:"id"`
	ReactionCount        map[string]int `json:"reaction_count"` // key = reaction emoji
	CurrentUserReactions []string       `json:"current_user_reactions"`
	CommentCount         *int64         `json:"comment_count"`
}

type TypingWS struct {
	SenderId int64 `json:"sender_id"`
}

type UserUpdateWS struct {
	Id     int64  `json:"id"`
	Status string `json:"status"` // online, offline
}

type InnerPacket struct {
	Payload      any
	HighPriority bool
}

// type CommentUpdateWS struct {
// 	Id                   int64          `json:"id"`
// 	ReactionCount        map[string]int `json:"reaction_count"`
// 	CurrentUserReactions []string       `json:"current_user_reaction"`
// }

// type UpdateMessageWS struct {
// 	Id                   int64          `json:"id"`
// 	ReactionCount        map[string]int `json:"reaction_count"`
// 	CurrentUserReactions []string       `json:"current_user_reactions"`
// }

// type PostUpdateWS struct {
// 	Id                   int64          `json:"id"`
// 	ReactionCount        map[string]int `json:"reaction_count"` // key = reaction emoji
// 	CurrentUserReactions []string       `json:"current_user_reactions"`
// 	CommentCount         int            `json:"comment_count"`
// }
