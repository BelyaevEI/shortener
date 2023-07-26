package route

import (
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/logger"
	m "github.com/BelyaevEI/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func New(h handlers.Handlers, log *logger.Logger) *chi.Mux {

	r := chi.NewRouter()

	// Подключаем middleware
	r.Use(m.Gzip)
	r.Use(log.Logger)
	r.Use(m.Cookie)

	r.Get("/{id}", h.ReplaceGET)
	r.Post("/api/shorten", h.PostAPI)
	r.Post("/", h.ReplacePOST)
	r.Get("/ping", h.PingDB)
	r.Post("/api/shorten/batch", h.PostAPIBatch)

	return r
}
