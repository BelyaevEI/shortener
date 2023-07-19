package cachestorage

import (
	"errors"

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

func (c *chache) GetShortUrl(inputURL string) (string, error) {
	if foundurl, ok := c.storageShortURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", errors.New("ErrNoRows")
}

func (c *chache) GetOriginUrl(inputURL string) (string, error) {
	if foundurl, ok := c.storageOriginURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", errors.New("ErrNoRows")
}

func (c *chache) Save(url1, url2 string) error {
	c.storageShortURL[url1] = url2
	c.storageOriginURL[url2] = url1
	return nil
}

func (c *chache) Ping() error {
	c.log.Log.Info("Work with file: no implement method Ping")
	return nil
}
