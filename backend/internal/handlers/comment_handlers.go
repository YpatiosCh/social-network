package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/eventSystem"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	urlparams "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/url_params"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) CreateCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NewCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		senderId := claims.UserId

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		id, createdAt, count, err := h.Db.CreateComment(ctx, req.Body, req.ParentId, senderId)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read comment")
			return
		}

		response := models.NewCommentResponse{
			Id:        id,
			CreatedAt: createdAt,
		}

		err = utils.WriteJSON(w, http.StatusCreated, response)
		if err != nil {
			fmt.Println("CreateComment failed to send response: ", err)
			return
		}

		comment, err := h.Db.GetCommentById(ctx, id, claims.UserId)
		if err != nil {
			fmt.Println("Failed to retrieve new comment: ", err)
			return
		}

		//comment count
		commentCount := models.ContentStatusUpdateWS{
			Id:           req.ParentId,
			CommentCount: &count,
		}

		countMessage := models.WSMessage{
			Type:    "status_update",
			Payload: commentCount,
		}

		err = eventSystem.SendEvent(fmt.Sprintf("su:%d", req.ParentId), countMessage, false)
		if err != nil {
			fmt.Println("CreateComment: Failed to send count event:", err.Error())
		}

		//actual comment
		commentMessage := models.WSMessage{
			Type:    "new_comment",
			Payload: comment,
		}

		err = eventSystem.SendEvent(fmt.Sprintf("nc:%d", req.ParentId), commentMessage, false)
		if err != nil {
			fmt.Println("CreateComment: Failed to send comment event:", err.Error())
		}
	}
}

func (h *Handlers) GetCommentsByPostIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /comments?parent_id=x&last_id=x&range=10
		UrlParams, err := urlparams.ParseUrlParams(r,
			urlparams.Params{"last_id": "int64", "range": "int", "parent_id": "int64"}, nil)
		fmt.Printf("endpoint %v with parsed params: %v\n", r.URL, UrlParams)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		request := models.CommentsRequest{}
		request.LastId, _ = UrlParams["last_id"].(int64)
		request.ParentID, _ = UrlParams["parent_id"].(int64)
		request.Range, _ = UrlParams["range"].(int)

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		userID := int64(claims.UserId)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		comments, haveMore, err := h.Db.GetCommentsByPostID(ctx, userID, request.ParentID, request.LastId, request.Range)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusNotFound, "no comments")
			return
		}

		response := models.CommentsResponse{
			Comments: comments,
			HaveMore: haveMore,
		}

		if len(response.Comments) > 0 {
			fmt.Println("total", response.Comments[0].ReactionCount, "user", response.Comments[0].CurrentUserReactions)
		}
		err = utils.WriteJSON(w, http.StatusOK, response)
		if err != nil {
			fmt.Println("GetCommentByPostId failed to send response: ", err)
			return
		}
	}
}
