package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/compres"
	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func RunServer() error {

	config.ParseFlags()

	r := chi.NewRouter()

	r.Get("/{id}", logger.WithLogging(compres.GzipMiddleware(handlers.ReplaceGET())))
	r.Post("/api/shorten", logger.WithLogging(compres.GzipMiddleware(handlers.PostAPI())))
	r.Post("/", logger.WithLogging(compres.GzipMiddleware(handlers.ReplacePOST())))

	return http.ListenAndServe(config.FlagRunAddr, r)
}
