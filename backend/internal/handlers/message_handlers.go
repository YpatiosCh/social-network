package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/eventSystem"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	urlparams "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/url_params"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) CreateMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NewMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		senderId := claims.UserId
		fmt.Println("senderID is", senderId)

		if strings.TrimSpace(req.Body) == "" || req.ReceiverId == 0 {
			utils.ErrorJSON(w, http.StatusBadRequest, "conversation_id, text and receiver_id are required")
			return
		}

		if req.ConversationId != 0 && req.ReceiverId == senderId {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid receiver")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		convID, messageID, sentAt, err := h.Db.CreateMessage(ctx, req.ConversationId, senderId, req.ReceiverId, req.Body)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not create message")
			return
		}

		resp := models.NewMessageResponse{
			ConversationId: convID,
			MessageId:      messageID,
			CreatedAt:      sentAt,
		}

		err = utils.WriteJSON(w, http.StatusCreated, resp)
		if err != nil {
			fmt.Println("CreateReaction failed to send response: ", err)
			return
		}

		//message to sender (resend the whole Feed(one user) with new count)
		//only send if client has multiple connections (other tabs, devices)
		if clients.CheckMulttipleConnections(senderId) {
			receiverInfo, err := h.Db.GetUserFeedById(ctx, senderId, req.ReceiverId)
			if err != nil {
				fmt.Println("CreateMessage: couldn't fetch receiver info: ", err)
				return
			}

			//check online status
			clients.UpdateUserStatus(&receiverInfo)

			eventData := models.WSMessage{
				Type:    "new_dm_count",
				Payload: receiverInfo,
			}

			if err := clients.SendDataToAllClientConns(senderId, eventData, true); err != nil {
				fmt.Println("CreateMessage: failed to publish receiver info:", err)
			}
		}

		//message to recipient (resend the whole Feed(one user) with new count)
		senderInfo, err := h.Db.GetUserFeedById(ctx, req.ReceiverId, senderId)
		if err != nil {
			fmt.Println("CreateMessage: couldn't fetch user: ", err)
			return
		}
		senderInfo.User.Status = "online" //assume online if submitting message

		eventData := models.WSMessage{
			Type:    "new_dm_count",
			Payload: senderInfo,
		}

		if err := clients.SendDataToAllClientConns(req.ReceiverId, eventData, true); err != nil {
			fmt.Println("CreateMessage: failed to publish User:", err)
		}

		//message to both conversation memmbers
		eventKey := eventSystem.CreateDmEventKey(senderId, req.ReceiverId)
		fmt.Println("Sending", eventKey)

		dm := models.Message{
			Id:                   messageID,
			ConversationId:       convID,
			SenderId:             senderId,
			Body:                 req.Body,
			CreatedAt:            &sentAt,
			CurrentUserReactions: []string{},
		}

		eventData = models.WSMessage{
			Type:    "new_dm",
			Payload: dm,
		}

		if err := eventSystem.SendEvent(eventKey, eventData, true); err != nil {
			fmt.Println("CreateMessage: failed to publish:", err)
		}
	}
}

func (h *Handlers) GetMessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		userId := claims.UserId
		urlParams, err := urlparams.ParseUrlParams(r, urlparams.Params{
			"id":     "int64",
			"before": "int64",
			"after":  "int64",
		}, urlparams.Params{"last_id": "int64"})
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid query")
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		req := models.MessagesDbRequest{
			Ctx:    ctx,
			UserId: userId,
			MessagesRequest: models.MessagesRequest{
				ConversationId: urlParams["id"].(int64),
				LastId:         urlParams["last_id"].(int64),
				Before:         urlParams["before"].(int64),
				After:          urlParams["after"].(int64),
			},
		}
		res, err := h.Db.GetConversationDms(req)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusNotFound, "conversation not found")
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, res)
		if err != nil {
			fmt.Println(err)
			fmt.Println("GetCommentByPostId failed to send response: ", err)
			return
		}
	}
}
