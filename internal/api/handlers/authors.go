package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/igferreira/quotes-api/internal/api"
	"github.com/igferreira/quotes-api/internal/repository"
	"github.com/igferreira/quotes-api/internal/service"
	"github.com/rs/zerolog/log"
)

// AuthorHandler handles author-related requests
type AuthorHandler struct {
	service *service.Service
}

// NewAuthorHandler creates a new author handler
func NewAuthorHandler(service *service.Service) *AuthorHandler {
	return &AuthorHandler{
		service: service,
	}
}

// Create handles POST /authors
func (h *AuthorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params repository.CreateAuthorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_REQUEST_BODY")
		return
	}

	// Validate input
	if params.Name == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("name is required"), "VALIDATION_ERROR")
		return
	}

	author, err := h.service.CreateAuthor(r.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create author")
		api.RespondError(w, http.StatusInternalServerError, err, "CREATE_AUTHOR_ERROR")
		return
	}

	api.RespondJSON(w, http.StatusCreated, author)
}

// GetByID handles GET /authors/{id}
func (h *AuthorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	author, err := h.service.GetAuthor(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to get author")
		api.RespondError(w, http.StatusNotFound, err, "AUTHOR_NOT_FOUND")
		return
	}

	api.RespondJSON(w, http.StatusOK, author)
}

// List handles GET /authors
func (h *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parsePaginationParams(r)

	authors, total, err := h.service.ListAuthors(r.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("failed to list authors")
		api.RespondError(w, http.StatusInternalServerError, err, "LIST_AUTHORS_ERROR")
		return
	}

	api.RespondPaginated(w, authors, total, params.Limit, params.Offset)
}

// Update handles PUT /authors/{id}
func (h *AuthorHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	var params repository.UpdateAuthorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_REQUEST_BODY")
		return
	}

	// Validate input
	if params.Name == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("name is required"), "VALIDATION_ERROR")
		return
	}

	author, err := h.service.UpdateAuthor(r.Context(), id, params)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to update author")
		api.RespondError(w, http.StatusInternalServerError, err, "UPDATE_AUTHOR_ERROR")
		return
	}

	api.RespondJSON(w, http.StatusOK, author)
}

// Delete handles DELETE /authors/{id}
func (h *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	err = h.service.DeleteAuthor(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to delete author")
		api.RespondError(w, http.StatusBadRequest, err, "DELETE_AUTHOR_ERROR")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Search handles GET /authors/search
func (h *AuthorHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("search query is required"), "VALIDATION_ERROR")
		return
	}

	params := parsePaginationParams(r)

	authors, total, err := h.service.SearchAuthors(r.Context(), query, params)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("failed to search authors")
		api.RespondError(w, http.StatusInternalServerError, err, "SEARCH_AUTHORS_ERROR")
		return
	}

	api.RespondPaginated(w, authors, total, params.Limit, params.Offset)
}
