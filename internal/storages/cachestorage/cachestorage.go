package cachestorage

import (
	"github.com/BelyaevEI/shortener/internal/logger"
)

type chache struct {
	storageShortURL  map[string]string
	storageOriginURL map[string]string
	log              *logger.Logger
}

func New(log *logger.Logger) *chache {
	return &chache{
		storageShortURL:  make(map[string]string),
		storageOriginURL: make(map[string]string),
		log:              log,
	}
}

func (c *chache) GetShortURL(inputURL string) (string, error) {
	if foundurl, ok := c.storageShortURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", nil
}

func (c *chache) GetOriginURL(inputURL string) (string, error) {
	if foundurl, ok := c.storageOriginURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", nil
}

func (c *chache) Save(shortURL, longURL string) error {
	c.storageShortURL[longURL] = shortURL
	c.storageOriginURL[shortURL] = longURL
	return nil
}

func (c *chache) Ping() error {
	c.log.Log.Info("Work with file: no implement method Ping")
	return nil
}
