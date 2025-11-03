package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/eventSystem"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// TODO doesn't work, should it?
func (h *Handlers) UsersReactedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UsersReactedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		// claims, ok := utils.GetValue[jwt.Claims](r, utils.ClaimsKey)
		// if !ok {panic(1)}
		// senderId := claims.UserId

		// ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		// defer cancel()
		// id, err := h.Db.GetReactions()
		// if err != nil {
		// 	utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read UsersReacted")
		// 	return
		// }

		response := models.UsersReactedResponse{
			Users: nil,
		}

		err := utils.WriteJSON(w, http.StatusCreated, response)
		if err != nil {
			fmt.Println("UsersReacted failed to send response: ", err)
			return
		}
	}
}

func (h *Handlers) CreateReactionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		userID := int64(claims.UserId)

		var in models.NewReactionRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&in); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "bad request: invalid JSON")
			return
		}
		if in.ContentId <= 0 {
			utils.ErrorJSON(w, http.StatusBadRequest, "bad request: invalid content_id")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		// create or delete
		//var total int64
		if in.New {
			if in.ReactionType == "" {
				utils.ErrorJSON(w, http.StatusBadRequest, "missing reaction type")
				return
			}
			_, _, err := h.Db.CreateReaction(ctx, userID, in.ContentId, in.ReactionType)
			if err != nil {
				utils.ErrorJSON(w, http.StatusInternalServerError, "failed to create reaction")
				return
			}
			//total = t
		} else {
			err := h.Db.DeleteReaction(ctx, userID, in.ContentId, in.ReactionType)
			if err != nil {
				fmt.Println(err)
				utils.ErrorJSON(w, http.StatusInternalServerError, "failed to delete reaction")
				return
			}
			//total = t
		}
		if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
			fmt.Println("Reaction handler: Failed to send json: ", err)
		}
		// HTTP response (no per-type counts)
		// users, err := h.Db.ListUsersReacted(ctx, in.ContentId)
		// if err != nil {
		// 	utils.ErrorJSON(w, http.StatusInternalServerError, "failed to load users")
		// 	return
		// }
		// resp := models.UsersReactedResponse{
		// 	Users: users,
		// 	Total: total,
		// }
		// _ = utils.WriteJSON(w, http.StatusOK, resp)

		//WS update for current user only
		counts, _, err := h.Db.ReactionCountsByType(ctx, in.ContentId)
		if err != nil {
			return
		}

		currentUserReactions, err := h.Db.CurrentUserReactions(ctx, userID, in.ContentId)
		if err == nil {
			payload := models.ContentStatusUpdateWS{
				Id:                   in.ContentId,
				CurrentUserReactions: currentUserReactions,
				ReactionCount:        counts,
			}

			if payload.CurrentUserReactions == nil {
				payload.CurrentUserReactions = []string{}
			}

			ev := models.WSMessage{
				Type:    "status_update",
				Payload: payload,
			}

			if err := clients.SendDataToAllClientConns(userID, ev, false); err != nil {
				fmt.Println("SetReaction: publish:", err)
			}
			fmt.Println("sent to current user")
		}

		// WS event for everyone subscribed
		payloadForAll := models.ContentStatusUpdateWS{
			Id:            in.ContentId,
			ReactionCount: counts, // map[string]int
		}

		evForAll := models.WSMessage{
			Type:    "status_update",
			Payload: payloadForAll,
		}

		if err := eventSystem.SendEvent(fmt.Sprintf("su:%d", in.ContentId), evForAll, false); err != nil {
			fmt.Println("SetReaction: publish:", err)
		}
		fmt.Println("sent to all users")
	}
}
