package cachestorage

import (
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

func (c *cache) GetShortURL(inputURL string) (string, error) {
	if foundurl, ok := c.storageShortURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", nil
}

func (c *cache) GetOriginURL(inputURL string) (string, error) {
	c.log.Log.Info(c.storageOriginURL)
	if foundurl, ok := c.storageOriginURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", nil
}

func (c *cache) Save(shortURL, longURL string) error {
	c.storageShortURL[longURL] = shortURL
	c.storageOriginURL[shortURL] = longURL
	return nil
}

func (c *cache) Ping() error {
	c.log.Log.Info("Work with internal storage: no implement method Ping")
	return nil
}
