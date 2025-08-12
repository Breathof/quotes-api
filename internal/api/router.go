package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/igferreira/quotes-api/internal/api/handlers"
	mw "github.com/igferreira/quotes-api/internal/api/middleware"
	"github.com/igferreira/quotes-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRouter creates a new router with all routes configured
func NewRouter(service *service.Service, db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.CORS)
	r.Use(middleware.Timeout(60)) // 60 second timeout

	// Health checks
	healthHandler := handlers.NewHealthHandler(db)
	r.Get("/healthz", healthHandler.Liveness)
	r.Get("/readyz", healthHandler.Readiness)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authors
		authorHandler := handlers.NewAuthorHandler(service)
		r.Route("/authors", func(r chi.Router) {
			r.Get("/", authorHandler.List)
			r.Post("/", authorHandler.Create)
			r.Get("/search", authorHandler.Search)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", authorHandler.GetByID)
				r.Put("/", authorHandler.Update)
				r.Delete("/", authorHandler.Delete)
			})
		})

		// Quotes
		quoteHandler := handlers.NewQuoteHandler(service)
		r.Route("/quotes", func(r chi.Router) {
			r.Get("/", quoteHandler.List)
			r.Post("/", quoteHandler.Create)
			r.Get("/search", quoteHandler.Search)
			r.Get("/random", quoteHandler.GetRandom)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", quoteHandler.GetByID)
				r.Put("/", quoteHandler.Update)
				r.Delete("/", quoteHandler.Delete)
			})
		})
	})

	return r
}
