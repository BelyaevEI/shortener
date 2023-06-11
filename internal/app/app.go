package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RunServer() error {

	config.ParseFlags()

	r := chi.NewRouter()
	r.Get("/{id}", handlers.ReplaceGET)
	r.Post("/", handlers.ReplacePOST)

	return http.ListenAndServe(config.FlagRunAddr, r)
}
