package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	requestcache "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/request_cache"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	urlparams "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/url_params"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

var reqCache = requestcache.New(1000)

func (h *Handlers) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NewPostRequest
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		senderId := claims.UserId

		// Parse the header for requestSignature and see if the req exists
		requestId, timeStamp := utils.ParseReqSignature(r, w) // not expecting anything on actionDetails on this endpoint
		if requestId == "" || timeStamp == "" {
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request: "+err.Error())
			return
		}

		categoryIds, err := categoryString2Slice(req.Categories)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, fmt.Sprintf("Bad request: bad value in categories: %s", err.Error()))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		id, createdAt, err := h.Db.CreatePost(ctx, req.Title, req.Body, senderId, categoryIds)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not create post")
			return
		}

		response := models.NewPostResponse{
			PostId:    id,
			CreatedAt: createdAt,
		}

		err = utils.WriteJSON(w, http.StatusCreated, response)
		if err != nil {
			fmt.Println("CreatePost failed to send response: ", err)
			return
		}

		post, err := h.Db.GetPostById(ctx, id, claims.UserId)
		if err != nil {
			fmt.Println("Failed to retrieve new post: ", err)
			return
		}

		//prepare message for clients
		message := models.WSMessage{
			Type:    "new_post",
			Payload: post,
		}

		err = clients.SendToAllClients(message, false)
		if err != nil {
			fmt.Println("CreatePost: Failed to send event:", err.Error())
		}

		// Add the request to the cache after successful processing
		reqCache.AddReq(requestcache.RequestMeta{
			Timestamp: timeStamp,
			UserId:    senderId,
		}, requestId)
	}
}

func (h *Handlers) GetPostByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.URL.Query().Get("post_id")
		fmt.Println(postIDStr)
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		userID := claims.UserId
		if strings.TrimSpace(postIDStr) == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "missing post_id")
			return
		}

		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil || postID <= 0 {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid post_id")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		post, err := h.Db.GetPostById(ctx, postID, userID)
		if err != nil {
			utils.ErrorJSON(w, http.StatusNotFound, "post not found")
			return
		}

		response := models.SinglePostResponse{
			Post: post,
		}

		err = utils.WriteJSON(w, http.StatusCreated, response)
		if err != nil {
			fmt.Println("CreatePost failed to send response: ", err)
			return
		}
	}
}

func (h *Handlers) GetPostsHandler() http.HandlerFunc {
	fmt.Println("getPostsHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}

		UrlParams, err := urlparams.ParseUrlParams(r,
			urlparams.Params{"last_id": "int64", "range": "int"},
			urlparams.Params{"categories": "int64 slice"})
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		req := models.PostsFeedRequest{
			LastSeenId: UrlParams["last_id"].(int64),
			Range:      UrlParams["range"].(int),
			Categories: UrlParams["categories"].([]int64),
		}

		if len(req.Categories) == 1 && req.Categories[0] == int64(0) {
			req.Categories = nil
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		res, err := h.Db.GetPosts(models.PostsFeedDbRequest{
			Ctx:              ctx,
			UserId:           claims.UserId,
			PostsFeedRequest: req,
		})
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusNotFound, "db error")
		}
		// fmt.Println(res)

		err = utils.WriteJSON(w, http.StatusOK, res)
		if err != nil {
			fmt.Println("CreatePost failed to send response: ", err)
			return
		}
	}
}

func categoryString2Slice(in []string) ([]int, error) {
	// Support both ["1","2"] and ["1,2"]
	flat := make([]string, 0, len(in))
	for _, s := range in {
		if strings.Contains(s, ",") {
			flat = append(flat, strings.Split(s, ",")...)
		} else {
			flat = append(flat, s)
		}
	}
	out := make([]int, 0, len(flat))
	for _, s := range flat {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		n, err := strconv.Atoi(s)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("bad cat id value: %q", s)
		}
		out = append(out, n)
	}
	return out, nil
}
