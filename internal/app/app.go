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
	flagRunAddr string
	chi         *chi.Mux
}

func RunServer() error {

	//Инициализируем сервис
	app := NewApp()
	return http.ListenAndServe(app.flagRunAddr, app.chi)
}

func NewApp() *App {

	// Парсинг переменных окружения
	cfg := config.ParseFlags()

	// Инициализируем хранилище
	s := storage.Init(cfg.FileStoragePath)

	// Создаем обьект handle
	h := handlers.New(cfg.ShortURL, s)

	// Создаем route
	r := route.New(h)

	return &App{
		flagRunAddr: cfg.FlagRunAddr,
		chi:         r,
	}
}
