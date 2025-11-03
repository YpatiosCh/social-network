package handlers

import (
	"context"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/database"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/middleware"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// TODO new context timeout middleware?

// /user
// /categories
// /posts?last_id=x&range=10&categories=x // category ids sep by comma "1,2,3"
// /users_feed?last_id=x&range=10&with_conv=x
// /post?id=x
// /comments?parent_id=x&last_id=x&range=10
// /dm?last_id=x&before=10&after=10 // including last_id
// /usersreacted?content_id=x&type=x&last_id=x&range=10 (OPTIONAL)

// /mark-seen, { method: 'POST',  body: "dm_ids": int}
// /new_post  {method: 'POST', body: { title: string, body: string, categories: string } } // category ids sep by comma "1,2,3"
// /new_comment {method: 'POST', body: {parent_id: int64, body: string}}
// /new_dm {
//     method: 'POST',
//     body: { conversation_id: int64, body: string }
//     }
// /new_reaction {
//      method: "POST",
//      body: { content_id: int64, type: string }
// }
// /typing {
//     method: "POST",
//     body: { recipient_id: int64 }
//     }

type Handlers struct {
	Db Database
}

type Database interface {
	CreateMessage(context.Context, int64, int64, int64, string) (int64, int64, time.Time, error)
	CreateUser(models.RegisterDbRequest) (models.RegisterDbResponse, error)
	CreatePost(context.Context, string, string, int64, []int) (int64, time.Time, error)
	CreateReaction(context.Context, int64, int64, string) (int64, int64, error)
	CreateComment(context.Context, string, int64, int64) (int64, time.Time, int64, error)
	DeleteReaction(context.Context, int64, int64, string) error
	UpdateUser(context.Context, int64, string, string, string, string, string, string, string) (int64, error)
	EditPost(context.Context, int64, string, string) (int64, error)
	EditComment(context.Context, int, string) (int, error)
	GetUserInfos(context.Context, int64) (*database.Users, error)
	GetPostById(context.Context, int64, int64) (*models.Post, error)
	GetCommentsByPostID(context.Context, int64, int64, int64, int) ([]models.Comment, bool, error)
	GetConversationDms(req models.MessagesDbRequest) (res models.MessagesDbResponse, err error)
	AddAuthUser(models.RegisterDbRequest) (models.RegisterDbResponse, error)
	Authenticate(context.Context, string, string) (int64, string, string, error)
	MarkSeen(models.MarkSeenDbRequest) (int64, error)
	AllCategories(context.Context) ([]models.Category, error)
	CategoriesById(context.Context, int64) (*models.Category, error)
	GetPosts(models.PostsFeedDbRequest) (models.PostsFeedDbResponse, error)
	ListUsersReacted(context.Context, int64) ([]models.UserReacted, error)
	ReactionCountsByType(context.Context, int64) (map[string]int, int64, error)
	GetUsersFeed(models.UsersFeedDbRequest) (models.UsersFeedDbResponse, error)
	GetCommentById(context.Context, int64, int64) (*models.Comment, error)
	CurrentUserReactions(context.Context, int64, int64) ([]string, error)
	GetUserFeedById(context.Context, int64, int64) (models.Feed, error)
}

func (h *Handlers) SetHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	Chain := middleware.Chain // ;)

	//NO AUTH
	mux.HandleFunc("/login", Chain().AllowedMethod("POST").Finalize(h.loginHandler()))
	mux.HandleFunc("/register", Chain().AllowedMethod("POST").Finalize(h.registerHandler()))

	//AUTH
	mux.Handle("/mark-seen", Chain().AllowedMethod("POST").Auth().Finalize(h.MarkSeenHandler()))
	mux.Handle("/typing", Chain().AllowedMethod("POST").Auth().Finalize(h.TypingHandler()))
	mux.Handle("/new_reaction", Chain().AllowedMethod("POST").Auth().Finalize(h.CreateReactionHandler()))
	mux.Handle("/new_post", Chain().AllowedMethod("POST").Auth().BindReqMeta().Finalize(h.CreatePostHandler()))
	mux.Handle("/new_comment", Chain().AllowedMethod("POST").Auth().Finalize(h.CreateCommentHandler()))
	mux.Handle("/new_dm", Chain().AllowedMethod("POST").Auth().Finalize(h.CreateMessageHandler()))

	//GET
	mux.HandleFunc("/ws", Chain().AllowedMethod("GET").Auth().Finalize(h.startWebsocket()))
	mux.HandleFunc("/logout", Chain().AllowedMethod("GET").Finalize(h.logoutHandler()))
	mux.Handle("/users_feed", Chain().AllowedMethod("GET").Auth().Finalize(h.UsersFeedHandler()))
	mux.Handle("/usersreacted", Chain().AllowedMethod("GET").Auth().Finalize(h.UsersReactedHandler()))
	mux.Handle("/dm", Chain().AllowedMethod("GET").Auth().Finalize(h.GetMessagesHandler()))
	mux.Handle("/posts", Chain().AllowedMethod("GET").Auth().Finalize(h.GetPostsHandler()))
	mux.Handle("/categories", Chain().AllowedMethod("GET").Auth().Finalize(h.CategoriesHandler()))
	mux.Handle("/comments", Chain().AllowedMethod("GET").Auth().Finalize(h.GetCommentsByPostIdHandler()))
	mux.Handle("/post", Chain().AllowedMethod("GET").Auth().Finalize(h.GetPostByIdHandler()))
	mux.Handle("/user", Chain().AllowedMethod("GET").Auth().Finalize(h.UserInfoHandler()))
	return mux
}
