package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/utils"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/shared/models"
)

func (h *Handlers) CategoriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		categories, err := h.Db.AllCategories(ctx)
		if err != nil {
			fmt.Println(err)
			utils.ErrorJSON(w, http.StatusInternalServerError, "Can not read Categories")
			return
		}

		response := models.CategoriesResponse{
			Categories: categories,
		}

		err = utils.WriteJSON(w, http.StatusOK, response)
		if err != nil {
			fmt.Println("Categories failed to send response: ", err)
			return
		}
	}
}

func (h *Handlers) CategoryByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		catIDStr := r.URL.Query().Get("category_id")
		if strings.TrimSpace(catIDStr) == "" {
			utils.ErrorJSON(w, http.StatusBadRequest, "missing category_id")
			return
		}

		categoryID, err := strconv.ParseInt(catIDStr, 10, 64)
		if err != nil || categoryID <= 0 {
			utils.ErrorJSON(w, http.StatusBadRequest, "invalid post_id")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		category, err := h.Db.CategoriesById(ctx, categoryID)
		if err != nil {
			utils.ErrorJSON(w, http.StatusNotFound, "post not found")
			return
		}
		resp := models.SingleCategoryResponse{
			Category: category,
		}
		err = utils.WriteJSON(w, http.StatusCreated, resp)
		if err != nil {
			fmt.Println("Get category failed to send response: ", err)
			return
		}
	}
}
