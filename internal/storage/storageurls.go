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

func New() *Storage {

	fileps := config.FileStoragePath

	//Будем считать, что в тестах будет путь /tmp/filename
	if _, err := os.Stat(filepath.Dir(fileps)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(fileps), 0755)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
	// открываем файл для записи
	file, err := os.OpenFile(fileps, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Open file with error: %s", err)
	}

	return &Storage{
		file:    file,
		encoder: json.NewEncoder(file),
		reader:  bufio.NewReader(file),
	}
}

func (s *Storage) Close() error {
	// закрываем файл
	return s.file.Close()
}

func (s *Storage) WriteURL(urls *models.StorageURL) error {
	return s.encoder.Encode(&urls)
}

func (s *Storage) ReadAllURLS() []models.StorageURL {

	var (
		read       [][]byte
		storageURL []models.StorageURL
	)

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
	return storageURL
}
