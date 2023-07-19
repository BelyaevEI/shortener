package utils

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
)

func GenerateRandomString(length int) string {
	// Задаем символы, из которых будет состоять случайная строка
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Генерируем случайную строку
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(result)

}

// Поиск короткой ссылки по длинной
func TryFoundShortURL(longURL string, s []models.StorageURL) (url string) {

	for _, ur := range s {
		if ur.OriginalURL == longURL {
			url = ur.ShortURL
			return url
		}
	}
	return ""
}

// Поиск длинной ссылки по короткой
func TryFoundOrigURL(shortURL string, s []models.StorageURL) (url string) {
	for _, ur := range s {
		if ur.ShortURL == shortURL {
			url = ur.OriginalURL
			return url
		}
	}
	return ""
}

// Макрос
func Response(w http.ResponseWriter, key, value, url string, statuscode int) {
	w.Header().Set(key, value)
	w.WriteHeader(statuscode)
	w.Write([]byte(url))
}

func ReadFile(path string, logger *logger.Logger) []models.StorageURL {

	var (
		read       [][]byte
		storageURL []models.StorageURL
	)

	// открываем файл для записи
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Log.Error(err)
		return nil
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

	return storageURL
}
