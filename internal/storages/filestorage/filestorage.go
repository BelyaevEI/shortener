package filestorage

import (
	"context"
	"encoding/json"
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

func (s *filestorage) Save(ctx context.Context, url1, url2 string, userID uint32) error {

	var longShortURL models.StorageURL

	// открываем файл для записи
	file, err := os.OpenFile(s.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	longShortURL.UserID, longShortURL.OriginalURL, longShortURL.ShortURL = userID, url2, url1
	encoder := json.NewEncoder(file)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return encoder.Encode(&longShortURL)
	}
}

func (s *filestorage) GetShortURL(ctx context.Context, inputURL string) (string, error) {

	var (
		storageURL []models.StorageURL
		foundurl   string
	)

	storageURL = utils.ReadFile(s.FileStoragePath, s.log)

	foundurl = utils.TryFoundShortURL(inputURL, storageURL)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return foundurl, nil
	}
}

func (s *filestorage) GetOriginURL(ctx context.Context, inputURL string) (string, bool, error) {

	var (
		storageURL []models.StorageURL
		foundurl   string
	)

	storageURL = utils.ReadFile(s.FileStoragePath, s.log)

	foundurl = utils.TryFoundOrigURL(inputURL, storageURL)

	select {
	case <-ctx.Done():
		return "", false, ctx.Err()
	default:
		return foundurl, false, nil
	}

}

func (s *filestorage) Ping(ctx context.Context) error {
	s.log.Log.Info("Work with file: no implement method Ping")
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *filestorage) GetUrlsUser(ctx context.Context, userID uint32) ([]models.StorageURL, error) {

	storageURL := utils.ReadFile(s.FileStoragePath, s.log)
	userURLS, err := utils.TryFoundUserURLS(userID, storageURL)
	if err != nil {
		return nil, err
	}
	return userURLS, nil
}

func (s *filestorage) UpdateDeletedFlag(ctx context.Context, data []string, userID uint32) {
}
