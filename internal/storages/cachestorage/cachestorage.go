package cachestorage

import (
	"context"

	"github.com/BelyaevEI/shortener/internal/logger"
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

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return foundurl, nil
	}
}

func (c *cache) Save(ctx context.Context, shortURL, longURL string) error {
	c.storageShortURL[longURL] = shortURL
	c.storageOriginURL[shortURL] = longURL
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (c *cache) Ping(ctx context.Context) error {
	c.log.Log.Info("Work with internal storage: no implement method Ping")
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
