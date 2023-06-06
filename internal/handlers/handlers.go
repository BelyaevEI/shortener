package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Вынести в отдельный каталог?
var short2long = make(map[string]string) //Словарь для получения полного URL по короткому
var long2short = make(map[string]string) //Словарь для получения короткого URL по полному

func ReplacePOST(w http.ResponseWriter, r *http.Request) {

	//Считать из тела запроса строку URL
	longURL, err := io.ReadAll(r.Body)
	if err != nil || string(longURL) == " " {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Проверим наличие короткой ссылки по длинной, если ее нет
	//то сгенерируем и запишем в словарь
	if shortURL, ok := long2short[string(longURL)]; ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	} else {
		short := generateRandomString(8)
		shortURL = "http://localhost:8080/" + short

		long2short[string(longURL)] = shortURL
		short2long[short] = string(longURL)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

func ReplaceGET(w http.ResponseWriter, r *http.Request) {

	var id string

	//получим ID из запроса
	idLong := r.URL.Path[1:]

	if strings.ContainsRune(idLong, '/') {
		id = strings.Split(idLong, "/")[0]
		if id == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		id = idLong
		if id == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//проверим по ID ссылку
	if longURL, ok := short2long[id]; ok {
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(longURL))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func generateRandomString(length int) string {
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
