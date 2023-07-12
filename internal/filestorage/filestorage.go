package filestorage

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/utils"
)

type filestorage struct {
	FileStoragePath string
}

func New(path string) *filestorage {

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Fatalf("Error: %s", err)
			return nil
		}
	}
	return &filestorage{FileStoragePath: path}
}

func (s *filestorage) Save(url1, url2 string) error {

	var longShortURL models.StorageURL

	// открываем файл для записи
	file, err := os.OpenFile(s.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
		return nil
	}

	defer file.Close()

	longShortURL.OriginalURL, longShortURL.ShortURL = url2, url1
	encoder := json.NewEncoder(file)

	return encoder.Encode(&longShortURL)
}

func (s *filestorage) Get(inputURL string) string {

	var (
		read       [][]byte
		storageURL []models.StorageURL
	)

	// открываем файл для записи
	file, err := os.OpenFile(s.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Ошибка при открытии %s", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	// Чтение из файла
	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		read = append(read, data)
	}

	// преобразуем данные из JSON-представления в структуру
	for _, line := range read {
		urls := models.StorageURL{}
		err := json.Unmarshal(line, &urls)
		if err == nil {
			storageURL = append(storageURL, urls)
		}
	}

	if foundurl1 := utils.TryFoundOrigURL(inputURL, storageURL); foundurl1 == "" {
		if foundurl2 := utils.TryFoundShortURL(inputURL, storageURL); foundurl1 != "" {
			return foundurl2
		}
	} else {
		return foundurl1
	}
	return ""
}
