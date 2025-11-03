package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.ErrorJSON(w, http.StatusMethodNotAllowed, "use POST")
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req models.RegisterRequest
		if err := dec.Decode(&req); err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		// Basic validation
		var identifier string
		if req.Email != nil || *req.Email != "" {
			identifier = strings.TrimSpace(*req.Email)
		} else if req.Username != nil || *req.Username != "" {
			identifier = strings.TrimSpace(*req.Username)
		} else {
			utils.ErrorJSON(w, http.StatusBadRequest, "username or email are required")
			return
		}

		if req.Password == "" ||
			req.ConfirmPassword == "" ||
			req.FirstName == "" ||
			req.LastName == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "username and password are required")
			return
		}

		if req.Password != req.ConfirmPassword {
			utils.ErrorJSON(w, http.StatusBadRequest, "passwords do not match")
			return
		}

		passwordHash, err := security.HashPassword(req.Password)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "could not hash password")
			return
		}

		if req.Email != nil {
			trimmedEmail := strings.TrimSpace(*req.Email)
			req.Email = &trimmedEmail
		}

		if req.Username != nil {
			trimmedUsername := strings.TrimSpace(*req.Username)
			req.Username = &trimmedUsername
		}

		// TODO: Trim rest of values

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbReq := models.RegisterDbRequest{
			Ctx:          ctx,
			Username:     req.Username,
			Email:        req.Email,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Gender:       req.Gender,
			Age:          req.Age,
			Avatar:       identifier[:1],
			Identifier:   identifier,
			PasswordHash: passwordHash,
		}

		dbRes, err := h.Db.AddAuthUser(dbReq)
		if err != nil {
			fmt.Println("Error creating user:", err)
			errMsg := strings.ToLower(err.Error())

			if strings.Contains(errMsg, "username") && strings.Contains(errMsg, "duplicate") {
				utils.ErrorJSON(w, http.StatusConflict, "Username already taken")
				return
			}
			if strings.Contains(errMsg, "email") && strings.Contains(errMsg, "duplicate") {
				utils.ErrorJSON(w, http.StatusConflict, "Email already taken")
				return
			}

			utils.ErrorJSON(w, http.StatusInternalServerError, "Could not create user")
			return

		}

		now := time.Now().Unix()
		exp := time.Now().AddDate(0, 6, 0).Unix()

		claims := security.Claims{
			UserId: dbRes.Id,
			Iat:    now,
			Nbf:    now,
			Exp:    exp,
		}

		token, err := security.CreateToken(claims)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "token generation failed")
			return
		}

		response := models.RegisterResponse{
			User: dbRes.RegisterResponseUser,
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Path:     "/",
			Expires:  time.Unix(exp, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		w.Header().Set("Location", fmt.Sprintf("/users/%d", dbRes.Id))
		_ = utils.WriteJSON(w, http.StatusCreated, response)

		user := models.UserOnFeed{
			Id:       dbRes.Id,
			Username: dbRes.UserName,
			Avatar:   dbRes.Avatar,
			Status:   "online", // just registered, must be online
		}

		conversation := models.ConversationDetails{}

		feed := models.Feed{
			User:                user,
			ConversationDetails: conversation,
		}

		//prepare message for clients
		message := models.WSMessage{
			Type:    "new_user",
			Payload: feed,
		}

		b, err := json.Marshal(message)
		if err != nil {
			fmt.Println("register struct broken?")
			panic(1)
		}

		fmt.Println("Broadcasting JSON:", string(b))
		err = clients.SendToAllClients(message, false)
		if err != nil {
			fmt.Println("CreatePost: Failed to send event:", err.Error())
		}

	}
}
