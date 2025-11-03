package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/clients"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/security"
	urlparams "platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/url_params"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) UsersFeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetValue[security.Claims](r, utils.ClaimsKey)
		if !ok {
			panic(1)
		}
		userID := claims.UserId
		UrlParams, err := urlparams.ParseUrlParams(r, urlparams.Params{"last_id": "int64", "range": "int"}, nil)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, "bad request")
			return
		}
		lastID, ok := UrlParams["last_id"].(int64)
		if !ok {
			panic(1)
		}
		pageSize, ok := UrlParams["range"].(int)
		if !ok {
			panic(1)
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		dbReq := models.UsersFeedDbRequest{
			Ctx:    ctx,
			UserId: userID,
			UsersFeedRequest: models.UsersFeedRequest{
				LastId: lastID,
				Range:  pageSize,
			},
		}

		dbResp, err := h.Db.GetUsersFeed(dbReq)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read UsersFeed")
			return
		}

		clients.UpdateUserStatuses(&dbResp)

		if err := utils.WriteJSON(w, http.StatusOK, dbResp.UsersFeedResponse); err != nil {
			fmt.Println("UsersFeed write json:", err)
		}
	}
}
