package cachestorage

import (
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
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
	return "", models.NoData
}

func (c *chache) GetOriginURL(inputURL string) (string, error) {
	if foundurl, ok := c.storageOriginURL[inputURL]; ok {
		return foundurl, nil
	}
	return "", models.NoData
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
