package utils

import (
	"math/rand"
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

func TryFoundOrigURL(shortUrl string, s []models.StorageURL) (url string) {

	for _, ur := range s {
		if ur.ShortURL == shortUrl {
			url = ur.OriginalURL
		}
	}
	return url
}

func TryFoundShortURL(u []byte, s []models.StorageURL) (url string) {
	longUrl := string(u)
	for _, ur := range s {
		if ur.OriginalURL == longUrl {
			url = ur.ShortURL
		}
	}
	return url
}
