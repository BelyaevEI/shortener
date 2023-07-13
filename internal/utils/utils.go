package utils

import (
	"math/rand"
	"net/http"
	"time"

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

func TryFoundShortURL(longURL string, s []models.StorageURL) (url string) {

	for _, ur := range s {
		if ur.OriginalURL == longURL {
			url = ur.ShortURL
			return url
		}
	}
	return ""
}

func TryFoundOrigURL(shortURL string, s []models.StorageURL) (url string) {
	for _, ur := range s {
		if ur.ShortURL == shortURL {
			url = ur.OriginalURL
			return url
		}
	}
	return ""
}

func Response(w http.ResponseWriter, key, value, url string, statuscode int) {
	w.Header().Set(key, value)
	w.WriteHeader(statuscode)
	w.Write([]byte(url))
}
