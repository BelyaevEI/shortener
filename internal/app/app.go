package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func RunServer() error {

	config.ParseFlags()

	r := chi.NewRouter()
	r.Get("/{id}", logger.WithLogging(handlers.ReplaceGET()))
	r.Post("/", logger.WithLogging(handlers.ReplacePOST()))

	return http.ListenAndServe(config.FlagRunAddr, r)
}
