package storage

import (
	"github.com/BelyaevEI/shortener/internal/cachestorage"
	"github.com/BelyaevEI/shortener/internal/database"
	"github.com/BelyaevEI/shortener/internal/filestorage"
	"github.com/BelyaevEI/shortener/internal/models"
)

type Storage struct {
	storage models.Storage
}

func Init(filepath, dbpath string) *Storage {
	if dbpath != "" {
		return &Storage{storage: database.New(dbpath)}
	}
	if filepath == " " {
		return &Storage{storage: cachestorage.New()}
	}
	return &Storage{storage: filestorage.New(filepath)}
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
