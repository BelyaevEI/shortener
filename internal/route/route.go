package route

import (
	"github.com/BelyaevEI/shortener/internal/handlers"
	m "github.com/BelyaevEI/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func New(h handlers.Handlers) *chi.Mux {

	r := chi.NewRouter()

	//Укажем middleware
	r.Use(m.Gzip)
	r.Use(m.Logger)

	r.Get("/{id}", h.ReplaceGET)
	r.Post("/api/shorten", h.PostAPI)
	r.Post("/", h.ReplacePOST)

	return r
}
