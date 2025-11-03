package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// /mark-seen, { method: 'POST',  body: "dm_ids": int}
func (h *Handlers) MarkSeenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		req := models.MarkSeenDbRequest{
			Ctx:    ctx,
			UserId: claims.UserId,
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}
		fmt.Println("conv", req.ConversationId, "DM", req.DmId)

		otherMemberId, err := h.Db.MarkSeen(req)
		fmt.Println(otherMemberId)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read MarkSeen")
			return
		}
		if err = utils.WriteJSON(w, http.StatusCreated, ok); err != nil {
			fmt.Println("MarkSeen failed to send response: ", err)
			return
		}

		//send updated feed to receiver
		receiverInfo, err := h.Db.GetUserFeedById(ctx, claims.UserId, otherMemberId)
		if err != nil {
			fmt.Println("Mark seen: couldn't fetch receiver info: ", err)
			return
		}

		//check online status
		clients.UpdateUserStatus(&receiverInfo)

		eventData := models.WSMessage{
			Type:    "marked_seen",
			Payload: receiverInfo,
		}

		if err := clients.SendDataToAllClientConns(claims.UserId, eventData, true); err != nil {
			fmt.Println("Mark Seen: failed to publish receiver info:", err)
		}
	}

}
