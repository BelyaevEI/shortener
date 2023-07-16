package storage

import (
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/storages/cachestorage"
	"github.com/BelyaevEI/shortener/internal/storages/database"
	"github.com/BelyaevEI/shortener/internal/storages/filestorage"
)

type Storage struct {
	storage models.Storage
}

func Init(filepath, dbpath string, log *logger.Logger) *Storage {
	if dbpath != "" {
		return &Storage{storage: database.New(dbpath, log)}
	}
	if filepath == " " {
		return &Storage{storage: cachestorage.New()}
	}
	return &Storage{storage: filestorage.New(filepath, log)}
}

func (s *Storage) GetURL(inputURL string) string {
	return s.storage.Get(inputURL)
}

func (s *Storage) SaveURL(url1, url2 string) error {
	return s.storage.Save(url1, url2)
}

func (s *Storage) Ping() error {
	return s.storage.Ping()
}
