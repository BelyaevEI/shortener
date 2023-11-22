// Package app create a new application for shortener service
// Commad to run server this service:
//
//	err = RunServer( )
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/handlers"
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/route"
	"github.com/BelyaevEI/shortener/internal/storages/storage"
	"golang.org/x/crypto/acme/autocert"
)

// This structure contain addres to run server and object chi
// type App struct {
// 	flagRunAddr string
// 	chi         *chi.Mux
// }

// For run server
func RunServer() {

	srv, err := runSrv()
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	sig := <-sigint
	log.Printf("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}

}

// Create a new application for service
func runSrv() (*http.Server, error) {

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

	srv := &http.Server{
		Addr:    cfg.FlagRunAddr,
		Handler: r,
	}

	if cfg.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.FlagRunAddr),
		}

		srv.TLSConfig = manager.TLSConfig()
		return srv, srv.ListenAndServeTLS("", "")

	}

	return srv, http.ListenAndServe(cfg.FlagRunAddr, r)
	// return &App{
	// 	flagRunAddr: cfg.FlagRunAddr,
	// 	chi:         r,
	// }

}
