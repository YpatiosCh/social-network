package handlers

import (
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
)

func (h *Handlers) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logout called")
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		err := utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "logged out successfully",
		})

		if err != nil{
			fmt.Println("failed to send logout ACK: ",err.Error())
			utils.ErrorJSON(w, http.StatusInternalServerError, "failed to send logout ACK")
		}
	}
}
