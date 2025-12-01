package application

import "time"

type GenericReq struct {
	RequesterId int64
	EntityId    int64
}

type GenericPaginatedReq struct {
	RequesterId int64
	EntityId    int64 //nil or 0 in case of personalized feed, or different model?
	Limit       int
	Offset      int
}

type hasRightToView struct {
	RequesterId         int64
	ParentEntityId      int64
	RequesterFollowsIds []int64
	RequesterGroups     []int64
}

//-------------------------------------------
// Posts
//-------------------------------------------
type Post struct {
	PostId          int64
	Body            string
	CreatorId       int64
	GroupId         int64  //can be nil (or 0?)
	Audience        string //everyone, selected, group, followers
	CommentsCount   int
	ReactionsCount  int
	ImagesCount     int
	LastCommentedAt time.Time //can be nil
	CreatedAt       time.Time
	UpdatedAt       time.Time //can be nil
	LikedByUser     bool
	FirstImage      string //can be nil (or "")
}

type CreatePostReq struct {
	CreatorId   int64
	Body        string
	GroupId     *int64  //can be nil (or 0?)
	Audience    string  //everyone, selected, group, followers
	AudienceIds []int64 //if audience=selected, audience ids go here, otherwise empty
}

type EditPostContentReq struct {
	RequesterId int64
	PostId      int64
	NewBody     string
}

type EditPostAudienceReq struct {
	RequesterId int64
	PostId      int64
	Audience    string
	AudienceIds []int64 //if audience=selected, audience ids go here, otherwise empty
}

type GetUserPostsReq struct {
	CreatorId        int64
	CreatorFollowers []int64 //from user service
	RequesterId      int64
	Limit            int
	Offset           int
}

type GetPersonalizedFeedReq struct {
	RequesterId         int64
	RequesterFollowsIds []int64 //from user service
	Limit               int
	Offset              int
}

//-------------------------------------------
// Comments
//-------------------------------------------

type Comment struct {
	CommentId      int64
	ParentId       int64
	Body           string
	CreatorId      int64
	ReactionsCount int
	ImagesCount    int
	CreatedAt      time.Time
	UpdatedAt      time.Time //can be nil
	LikedByUser    bool
	FirstImage     string //can be nil (or "")
}

type CreateCommentReq struct {
	CreatorId int64
	ParentId  int64
	Body      string
}

type EditCommentReq struct {
	CreatorId int64
	CommentId int64
	Body      string
}

//-------------------------------------------
// Events
//-------------------------------------------

type Event struct {
	EventId       int64
	Title         string
	Body          string
	CreatorId     int64
	GroupId       int64
	EventDate     time.Time
	StillValid    bool //still not sure this is needed
	GoingCount    int
	NotGoingCount int
	ImagesCount   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateEventReq struct {
	Title     string
	Body      string
	CreatorId int64
	GroupId   int64
	EventDate time.Time
}

type EditEventReq struct {
	EventId     int64
	RequesterId int64
	Title       string
	Body        string
	EventDate   time.Time
}

type RespondToEventReq struct {
	EventId     int64
	ResponderId int64
	Going       bool
}

//-------------------------------------------
// Reactions
//-------------------------------------------

//-------------------------------------------
// Images
//-------------------------------------------

type InsertImagesReq struct {
	PostId int64
	Images []string
}
