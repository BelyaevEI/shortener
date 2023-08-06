package utils

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
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

// Генерация уникального ID для пользователя
func GenerateUniqueID() uint32 {

	time := time.Now().UnixNano()

	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Конвертируем случайное число в uint64
	randomNumber := binary.BigEndian.Uint32(randomBytes)

	// Добавляем к времени случайное число
	uniqueNumber := uint32(time) + randomNumber

	return uniqueNumber
}

// Поиск ссылок в файле по юзеру
func TryFoundUserURLS(userID uint32, s []models.StorageURL) ([]models.StorageURL, error) {
	store := make([]models.StorageURL, 0)

	for _, line := range s {
		if line.UserID == userID {
			store = append(store, line)
		}
	}

	if len(store) == 0 {
		return nil, nil
	}

	return store, nil
}

func RemoveDuplicate(deleteURLS []string) []models.DeleteURL {
	var result []models.DeleteURL

	dd := make(map[string]struct{})

	for _, v := range deleteURLS {
		if _, ok := dd[v]; !ok {
			dd[v] = struct{}{}

			var res models.DeleteURL
			res.ShortURL = v
			result = append(result, res)
		}
	}
	return result
}

func MarkDeletion(userURLS []models.StorageURL, deleteURLS []string) []models.DeleteURL {
	var del models.DeleteURL
	markDel := make([]models.DeleteURL, 0)

	for _, varDel := range deleteURLS {
		for _, v := range userURLS {
			if !v.DeletedFlag && v.ShortURL == varDel {
				del.ShortURL = v.ShortURL
				markDel = append(markDel, del)
			}
		}
	}
	return markDel
}

func Generator(doneCh chan struct{}, input []models.StorageURL) chan models.StorageURL {
	inputCh := make(chan models.StorageURL)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}
