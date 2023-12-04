// Storage use internal memory
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

// Create a new storage
func New(log *logger.Logger) *cache {
	return &cache{
		storageShortURL:  make(map[string]string),
		storageOriginURL: make(map[string]string),
		log:              log,
	}
}

// Find short url
func (c *cache) GetShortURL(ctx context.Context, inputURL string) (string, error) {
	foundurl := c.storageShortURL[inputURL]
	return foundurl, nil
}

// Find origin url
func (c *cache) GetOriginURL(ctx context.Context, inputURL string) (string, bool, error) {
	foundurl := c.storageOriginURL[inputURL]
	return foundurl, false, nil
}

// Save urls to memory
func (c *cache) Save(ctx context.Context, shortURL, longURL string, userID uint32) error {
	c.storageShortURL[longURL] = shortURL
	c.storageOriginURL[shortURL] = longURL
	return nil
}

// This is mock function
func (c *cache) Ping(ctx context.Context) error {
	c.log.Log.Info("Work with internal storage: no implement method Ping")
	return nil
}

// This is mock function
func (c *cache) GetUrlsUser(ctx context.Context, userID uint32) ([]models.StorageURL, error) {
	return nil, nil
}

// This is mock function
func (c *cache) UpdateDeletedFlag(data models.DeleteURL) {
}

func (c *cache) GetStatistic() models.Statistic {
	var stat models.Statistic
	return stat
}
