package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/igferreira/quotes-api/internal/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *pgxpool.Pool
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Liveness handles GET /healthz
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	api.RespondJSON(w, http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  map[string]string{},
	})
}

// Readiness handles GET /readyz
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	services := make(map[string]string)
	status := http.StatusOK
	overallStatus := "ok"

	// Check database
	if err := h.db.Ping(ctx); err != nil {
		services["database"] = "unhealthy: " + err.Error()
		status = http.StatusServiceUnavailable
		overallStatus = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	api.RespondJSON(w, status, HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
	})
}
