package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login handler called")
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req models.LoginRequest
		if err := dec.Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid json")
			return
		}

		if strings.TrimSpace(req.Identifier) == "" || req.Password == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "identifier and password are required")
			return
		}

		fmt.Println("login received: ", req)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		usrId, username, avatar, err := h.Db.Authenticate(ctx, req.Identifier, req.Password)
		if err != nil {
			fmt.Println("failed to authenticate:", err)
			utils.ErrorJSON(w, http.StatusUnauthorized, "Invalid username or password")
			return

		}

		now := time.Now().Unix()
		exp := time.Now().AddDate(0, 6, 0).Unix() // six months from now

		claims := security.Claims{
			UserId: int64(usrId),
			Iat:    now,
			Exp:    exp,
		}

		token, err := security.CreateToken(claims)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "token generation failed")
			return
		}

		response := models.LoginResponse{
			User: models.LoginResponseUser{Id: usrId, Username: username, Avatar: avatar},
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Path:     "/",
			Expires:  time.Unix(exp, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		err = utils.WriteJSON(w, http.StatusOK, response)
		if err != nil {
			fmt.Println("Login handler: Failed to send json: ", err)
		}
	}
}
