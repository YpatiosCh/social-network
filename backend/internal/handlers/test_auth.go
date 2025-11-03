package handlers

import (
	"net/http"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
)

func privateHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := utils.GetValue[security.Claims](r, utils.ClaimsKey)
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "private ok",
		"user_id": claims.UserId,
		"exp":     claims.Exp,
	})
}
