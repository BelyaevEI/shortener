package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/route"
	"github.com/BelyaevEI/shortener/internal/storages/storage"
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

	//Создаем логгер
	log := logger.New()

	// Парсинг переменных окружения
	cfg := config.ParseFlags()

	// Инициализируем хранилище
	s := storage.Init(cfg.FileStoragePath, cfg.DBpath, log)

	// Создаем обьект handle
	h := handlers.New(cfg.ShortURL, s, log)

	// Создаем route
	r := route.New(h, log)

	return &App{
		flagRunAddr: cfg.FlagRunAddr,
		chi:         r,
	}

}
