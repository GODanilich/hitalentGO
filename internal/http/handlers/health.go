package handlers

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "fail",
			"db":     "unavailable",
		})
		return
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status": "fail",
			"db":     "unavailable",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
