package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func Shortener(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handlers.ReplacePOST(w, r)
	case http.MethodGet:
		handlers.ReplaceGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RunServer() error {

	r := chi.NewRouter()
	r.Get("/{id}", handlers.ReplaceGET)
	r.Post("/", handlers.ReplacePOST)
	return http.ListenAndServe(":8080", r)
}
