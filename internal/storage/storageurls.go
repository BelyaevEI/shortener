package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/BelyaevEI/shortener/internal/models"
)

type Storage struct {
	file    *os.File
	encoder *json.Encoder
	reader  *bufio.Reader
}

func NewStorage() *Storage {

	// открываем файл для записи в конец
	file, err := os.OpenFile("short-url-db.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return &Storage{}
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
		read        [][]byte
		storageURLS []models.StorageURL
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
			storageURLS = append(storageURLS, urls)
		}
	}
	return storageURLS
}
