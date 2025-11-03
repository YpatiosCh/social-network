package handlers

// TODO (by alex: ) Don't know if this needs to exist

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

// // This is only for admin use to create users without authentication
// func  (h *Handlers) CreateUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		if r.Method != http.MethodPost {
// 			utils.ErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
// 			return
// 		}

// 		var req models.User
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
// 			return
// 		}
// 		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
// 		defer cancel()
// 		id, err := h.Db.CreateUser(ctx, pool, req.Username, req.Gender, req.FirstName, req.LastName, req.Email, req.Age, req.Avatar)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "could not create user")
// 			return
// 		}

// 		var sentAt time.Time
// 		if err := pool.QueryRow(ctx, `SELECT created_at FROM users WHERE id = $1`, id).Scan(&sentAt); err != nil {
// 			// fallback if column has default and query failed
// 			sentAt = time.Now()
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(models.CreateResponse{
// 			Success: true,
// 			ID:      id,
// 			SentAt:  sentAt,
// 		})
// 	}
// }

// func  (h *Handlers) UpdateUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		if r.Method != http.MethodPut {
// 			utils.ErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
// 			return
// 		}
// 		var upd models.User
// 		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
// 			utils.ErrorJSON(w, http.StatusBadRequest, "Error decoding users request")
// 			return
// 		}
// 		if upd.Id == 0 {
// 			utils.ErrorJSON(w, http.StatusBadRequest, " User id required")
// 		}
// 		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
// 		defer cancel()
// 		id, err := h.Db.UpdateUser(ctx, pool, upd.Id,
// 			upd.Username, upd.Gender, upd.FirstName, upd.LastName,
// 			upd.Email, upd.Age, upd.Avatar)
// 		if err != nil {
// 			utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update user")
// 			return
// 		}
// 		var updatedAt *time.Time
// 		var nt sql.NullTime
// 		if err := pool.QueryRow(ctx, `SELECT updated_at FROM users WHERE id = $1`, id).Scan(&nt); err == nil && nt.Valid {
// 			t := nt.Time
// 			updatedAt = &t
// 		}

// 		_ = json.NewEncoder(w).Encode(models.UpdateResponse{
// 			Success:   true,
// 			ID:        id,
// 			UpdatedAt: updatedAt,
// 		})
// 	}
// }

// // func (h *Handlers)  CreateAuthUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
// // 	return func(w http.ResponseWriter, r *http.Request) {
// // 		w.Header().Set("Content-Type", "application/json")
// // 		if r.Method != http.MethodPost {
// // 			utils.ErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
// // 			return
// // 		}

// // 		var req models.User
// // 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// // 			utils.ErrorJSON(w, http.StatusBadRequest, "Bad request")
// // 			return
// // 		}
// // 		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
// // 		defer cancel()
// // 		id, err := h.Db.AddAuthUser(ctx, pool, req.Username, req.Password, req.Gender, req.FirstName, req.LastName, req.Email, req.Age, req.Avatar)
// // 		if err != nil {
// // 			utils.ErrorJSON(w, http.StatusInternalServerError, "could not create user")
// // 			return
// // 		}

// // 		var sentAt time.Time
// // 		if err := pool.QueryRow(ctx, `SELECT created_at FROM users WHERE id = $1`, id).Scan(&sentAt); err != nil {
// // 			// fallback if column has default and query failed
// // 			sentAt = time.Now()
// // 		}

// // 		w.Header().Set("Content-Type", "application/json")
// // 		w.WriteHeader(http.StatusCreated)
// // 		json.NewEncoder(w).Encode(models.CreateResponse{
// // 			Success: true,
// // 			ID:      id,
// // 			SentAt:  sentAt,
// // 		})
// // 	}

// // }

func (h *Handlers) UserInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		senderId := claims.UserId

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		info, err := h.Db.GetUserInfos(ctx, senderId)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read UserInfo")
			return
		}

		response := models.UserResponse{
			Id:       info.Id,
			Username: info.Username,
			Avatar:   info.Avatar,
		}
		fmt.Println(response)

		err = utils.WriteJSON(w, http.StatusCreated, response)
		if err != nil {
			fmt.Println("UserInfo failed to send response: ", err)
			return
		}
	}
}
