package storage

import (
	"github.com/BelyaevEI/shortener/internal/cachestorage"
	"github.com/BelyaevEI/shortener/internal/filestorage"
	"github.com/BelyaevEI/shortener/internal/models"
)

type Storage struct {
	storage models.Storage
}

func Init(path string) *Storage {
	if path == " " {
		return &Storage{storage: cachestorage.New()}
	}
	return &Storage{storage: filestorage.New(path)}
}

func (s *Storage) GetURL(inputURL string) string {
	return s.storage.Get(inputURL)
}

func (s *Storage) SaveURL(url1, url2 string) error {
	return s.storage.Save(url1, url2)
}
