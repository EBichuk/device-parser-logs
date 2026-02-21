package controller

import (
	"context"
	"device-parser-logs/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
)

type service interface {
	GetDeviceLogs(context.Context, string, int, int) (*models.PaginationResult, error)
}

type Handler struct {
	s      service
	logger *slog.Logger
}

func NewHandler(s service, log *slog.Logger) *Handler {
	return &Handler{
		s:      s,
		logger: log,
	}
}

func (h *Handler) GetDeviceLogsByGuid(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "guid")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 5 {
		limit = 5
	}

	files, err := h.s.GetDeviceLogs(r.Context(), guid, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}
