package route

import (
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/midllewares"
	"github.com/go-chi/chi/v5"
)

func New(h handlers.Handlers) *chi.Mux {

	r := chi.NewRouter()

	//Укажем middleware
	r.Use(midllewares.Gzip)
	r.Use(midllewares.Logger)

	r.Get("/{id}", h.ReplaceGET)
	r.Post("/api/shorten", h.PostAPI)
	r.Post("/", h.ReplacePOST)

	return r
}
