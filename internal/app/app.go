// Package app create a new application for shortener service
// Commad to run server this service:
//
//	err = RunServer( )
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

// This structure contain addres to run server and object chi
type App struct {
	flagRunAddr string
	chi         *chi.Mux
}

// For run server
func RunServer() error {
	//Инициализируем сервис
	app := NewApp()
	return http.ListenAndServe(app.flagRunAddr, app.chi)

}

// Create a new application for service
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
