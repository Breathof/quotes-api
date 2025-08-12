package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Error().Err(err).Msg("failed to encode response")
		}
	}
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, status int, err error, code string) {
	RespondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: err.Error(),
		Code:    code,
	})
}

// RespondPaginated sends a paginated response
func RespondPaginated(w http.ResponseWriter, data interface{}, total int64, limit, offset int32) {
	RespondJSON(w, http.StatusOK, PaginatedResponse{
		Data: data,
		Meta: PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}
