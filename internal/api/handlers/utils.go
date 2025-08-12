package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/igferreira/quotes-api/internal/repository"
)

// Default pagination values
const (
	DefaultLimit  = 20
	DefaultOffset = 0
	MaxLimit      = 100
)

// ErrValidation creates a validation error
func ErrValidation(msg string) error {
	return errors.New(msg)
}

// parsePaginationParams parses pagination parameters from request
func parsePaginationParams(r *http.Request) repository.ListParams {
	limit := DefaultLimit
	offset := DefaultOffset

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > MaxLimit {
				limit = MaxLimit
			}
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return repository.ListParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
}
