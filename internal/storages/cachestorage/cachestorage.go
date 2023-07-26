package cachestorage

import (
	"context"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
)

type cache struct {
	storageShortURL  map[string]string
	storageOriginURL map[string]string
	log              *logger.Logger
}

func New(log *logger.Logger) *cache {
	return &cache{
		storageShortURL:  make(map[string]string),
		storageOriginURL: make(map[string]string),
		log:              log,
	}
}

func (c *cache) GetShortURL(ctx context.Context, inputURL string) (string, error) {
	foundurl := c.storageShortURL[inputURL]

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return foundurl, nil
	}

}

func (c *cache) GetOriginURL(ctx context.Context, inputURL string) (string, error) {
	foundurl := c.storageOriginURL[inputURL]

	return foundurl, nil

}

func (c *cache) Save(ctx context.Context, shortURL, longURL string, userID uint64) error {
	c.storageShortURL[longURL] = shortURL
	c.storageOriginURL[shortURL] = longURL
	return nil
}

func (c *cache) Ping(ctx context.Context) error {
	c.log.Log.Info("Work with internal storage: no implement method Ping")
	return nil
}

func (c *cache) GetUrlsUser(ctx context.Context, userID uint64) ([]models.StorageURL, error) {
	return nil, nil
}
