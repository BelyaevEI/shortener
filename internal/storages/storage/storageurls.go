package storage

import (
	"context"

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
		return &Storage{storage: cachestorage.New(log)}
	}
	return &Storage{storage: filestorage.New(filepath, log)}
}

func (s *Storage) GetOriginalURL(ctx context.Context, inputURL string) (string, bool, error) {
	return s.storage.GetOriginURL(ctx, inputURL)
}

func (s *Storage) GetShortenURL(ctx context.Context, inputURL string) (string, error) {
	return s.storage.GetShortURL(ctx, inputURL)
}

func (s *Storage) SaveURL(ctx context.Context, url1, url2 string, userID uint32) error {
	return s.storage.Save(ctx, url1, url2, userID)
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s *Storage) GetUrlsUser(ctx context.Context, userID uint32) ([]models.StorageURL, error) {
	return s.storage.GetUrlsUser(ctx, userID)
}

func (s *Storage) UpdateDeletedFlag(ctx context.Context, data []models.DeleteURL) error {
	return s.storage.UpdateDeletedFlag(ctx, data)
}
