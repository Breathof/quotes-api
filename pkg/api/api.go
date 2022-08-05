package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/quote", func(r chi.Router) {
		r.Get("/", getQuote)
		r.Post("/", addQuote)
		r.Put("/{quoteID}", updateQuote)
		r.Delete("/{quoteID}", deleteQuote)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			log.Printf("error processing request: %v\n", err)
		}
	})

	return r
}

func getQuote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get quote"))
}

func addQuote(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create quote"))
}

func updateQuote(w http.ResponseWriter, r *http.Request) {
	quoteID := chi.URLParam(r, "quoteID")
	w.Write([]byte("update quote " + quoteID))
}

func deleteQuote(w http.ResponseWriter, r *http.Request) {
	quoteID := chi.URLParam(r, "quoteID")
	w.Write([]byte("delete quote " + quoteID))
}
