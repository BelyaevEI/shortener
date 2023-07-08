package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/route"
	"github.com/BelyaevEI/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type App struct {
	cfg    config.Parameters
	handle handlers.Handlers
	chi    *chi.Mux
}

func RunServer() error {

	//Инициализируем сервис
	app := NewApp()
	return http.ListenAndServe(app.cfg.FlagRunAddr, app.chi)
}

func NewApp() *App {

	// Парсинг переменных окружения
	cfg := config.ParseFlags()

	//Работа с файлом
	s := storage.New(cfg)

	//Создаем обьект handle
	h := handlers.New(cfg, s)

	// Создаем route
	r := route.New(h)

	return &App{
		cfg:    cfg,
		handle: h,
		chi:    r,
	}
}
