package filestorage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/utils"
)

type filestorage struct {
	FileStoragePath string
	log             *logger.Logger
}

func New(path string, log *logger.Logger) *filestorage {

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Log.Error(err)
			return nil
		}
	}
	return &filestorage{
		FileStoragePath: path,
		log:             log}
}

func (s *filestorage) Save(url1, url2 string) error {

	var longShortURL models.StorageURL

	// открываем файл для записи
	file, err := os.OpenFile(s.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	longShortURL.OriginalURL, longShortURL.ShortURL = url2, url1
	encoder := json.NewEncoder(file)

	return encoder.Encode(&longShortURL)
}

func (s *filestorage) GetShortURL(inputURL string) (string, error) {

	var (
		storageURL []models.StorageURL
		foundurl   string
	)

	storageURL = utils.ReadFile(s.FileStoragePath, s.log)

	if foundurl := utils.TryFoundShortURL(inputURL, storageURL); foundurl != "" {
		return "", errors.New("ErrNoRows")
	}

	return foundurl, nil
}

func (s *filestorage) GetOriginURL(inputURL string) (string, error) {

	var (
		storageURL []models.StorageURL
		foundurl   string
	)

	storageURL = utils.ReadFile(s.FileStoragePath, s.log)

	if foundurl := utils.TryFoundOrigURL(inputURL, storageURL); foundurl != "" {
		return "", errors.New("ErrNoRows")
	}

	return foundurl, nil
}

func (s *filestorage) Ping() error {
	s.log.Log.Info("Work with file: no implement method Ping")
	return nil
}
