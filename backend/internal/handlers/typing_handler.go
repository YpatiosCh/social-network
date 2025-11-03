package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/eventSystem"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) TypingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.TypingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
			return
		}

		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		typerId := claims.UserId

		err := utils.WriteJSON(w, http.StatusOK, nil)
		if err != nil {
			fmt.Println("Typing failed to send response: ", err)
			return
		}

		eventKey := eventSystem.CreateDmEventKey(typerId, req.RecipientId)
		fmt.Println(eventKey)

		payload := models.TypingWS{
			SenderId: typerId,
		}

		eventData := models.WSMessage{
			Type:    "typing",
			Payload: payload,
		}

		if err := eventSystem.SendEvent(eventKey, eventData, true); err != nil {
			fmt.Println("Typing: failed to publish:", err)
		}
	}
}
