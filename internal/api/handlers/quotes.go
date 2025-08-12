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

// QuoteHandler handles quote-related requests
type QuoteHandler struct {
	service *service.Service
}

// NewQuoteHandler creates a new quote handler
func NewQuoteHandler(service *service.Service) *QuoteHandler {
	return &QuoteHandler{
		service: service,
	}
}

// Create handles POST /quotes
func (h *QuoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params repository.CreateQuoteParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_REQUEST_BODY")
		return
	}

	// Validate input
	if params.Content == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("content is required"), "VALIDATION_ERROR")
		return
	}
	if params.AuthorID <= 0 {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("valid author_id is required"), "VALIDATION_ERROR")
		return
	}

	quote, err := h.service.CreateQuote(r.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create quote")
		api.RespondError(w, http.StatusInternalServerError, err, "CREATE_QUOTE_ERROR")
		return
	}

	api.RespondJSON(w, http.StatusCreated, quote)
}

// GetByID handles GET /quotes/{id}
func (h *QuoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	quote, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to get quote")
		api.RespondError(w, http.StatusNotFound, err, "QUOTE_NOT_FOUND")
		return
	}

	api.RespondJSON(w, http.StatusOK, quote)
}

// List handles GET /quotes
func (h *QuoteHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parsePaginationParams(r)

	// Check if filtering by author
	authorIDStr := r.URL.Query().Get("author_id")
	if authorIDStr != "" {
		authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err, "INVALID_AUTHOR_ID")
			return
		}

		quotes, err := h.service.ListQuotesByAuthor(r.Context(), authorID, params)
		if err != nil {
			log.Error().Err(err).Int64("author_id", authorID).Msg("failed to list quotes by author")
			api.RespondError(w, http.StatusInternalServerError, err, "LIST_QUOTES_ERROR")
			return
		}

		api.RespondJSON(w, http.StatusOK, quotes)
		return
	}

	quotes, total, err := h.service.ListQuotes(r.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("failed to list quotes")
		api.RespondError(w, http.StatusInternalServerError, err, "LIST_QUOTES_ERROR")
		return
	}

	api.RespondPaginated(w, quotes, total, params.Limit, params.Offset)
}

// Update handles PUT /quotes/{id}
func (h *QuoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	var params repository.UpdateQuoteParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_REQUEST_BODY")
		return
	}

	// Validate input
	if params.Content == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("content is required"), "VALIDATION_ERROR")
		return
	}
	if params.AuthorID <= 0 {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("valid author_id is required"), "VALIDATION_ERROR")
		return
	}

	quote, err := h.service.UpdateQuote(r.Context(), id, params)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to update quote")
		api.RespondError(w, http.StatusInternalServerError, err, "UPDATE_QUOTE_ERROR")
		return
	}

	api.RespondJSON(w, http.StatusOK, quote)
}

// Delete handles DELETE /quotes/{id}
func (h *QuoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		api.RespondError(w, http.StatusBadRequest, err, "INVALID_ID")
		return
	}

	err = h.service.DeleteQuote(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to delete quote")
		api.RespondError(w, http.StatusInternalServerError, err, "DELETE_QUOTE_ERROR")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Search handles GET /quotes/search
func (h *QuoteHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		api.RespondError(w, http.StatusBadRequest, ErrValidation("search query is required"), "VALIDATION_ERROR")
		return
	}

	params := parsePaginationParams(r)

	quotes, total, err := h.service.SearchQuotes(r.Context(), query, params)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("failed to search quotes")
		api.RespondError(w, http.StatusInternalServerError, err, "SEARCH_QUOTES_ERROR")
		return
	}

	api.RespondPaginated(w, quotes, total, params.Limit, params.Offset)
}

// GetRandom handles GET /quotes/random
func (h *QuoteHandler) GetRandom(w http.ResponseWriter, r *http.Request) {
	quote, err := h.service.GetRandomQuote(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get random quote")
		api.RespondError(w, http.StatusInternalServerError, err, "GET_RANDOM_QUOTE_ERROR")
		return
	}

	api.RespondJSON(w, http.StatusOK, quote)
}
