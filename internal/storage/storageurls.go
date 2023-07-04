package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/models"
)

type Storage struct {
	file    *os.File
	encoder *json.Encoder
	reader  *bufio.Reader
}

var (
	Manage     Storage
	storageURL []models.StorageURL
)

func Init() {

	//Будем считать, что в тестах будет путь /tmp/filename
	if _, err := os.Stat(filepath.Dir(config.FileStoragePath)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(config.FileStoragePath), 0755)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	// открываем файл для записи
	file, err := os.OpenFile(config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Open file with error: %s", err)
	}

	Manage = Storage{
		file:    file,
		encoder: json.NewEncoder(file),
		reader:  bufio.NewReader(file),
	}

	// Прочитаем содержимое файла
	Manage.ReadAllURLS()
}

func (s *Storage) Close() error {
	// закрываем файл
	return s.file.Close()
}

func (s *Storage) WriteURL(urls *models.StorageURL) error {
	return s.encoder.Encode(&urls)
}

func (s *Storage) ReadAllURLS() {

	var read [][]byte

	// Чтение из файла
	for {
		data, err := s.reader.ReadBytes('\n')
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
}

func (s *Storage) TryFoundOrigURL(shortURL string) (url string) {
	for _, ur := range storageURL {
		if ur.ShortURL == shortURL {
			url = ur.OriginalURL
		}
	}
	return url
}

func (s *Storage) TryFoundShortURL(u []byte) (url string) {
	longURL := string(u)
	for _, ur := range storageURL {
		if ur.OriginalURL == longURL {
			url = ur.ShortURL
		}
	}
	return url
}
