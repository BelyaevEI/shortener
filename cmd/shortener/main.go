package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var short2long = make(map[string]string) //Словарь для получения полного URL по короткому
var long2short = make(map[string]string) //Словарь для получения короткого URL по полному

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortener)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func shortener(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		replacePOST(w, r)
	case http.MethodGet:
		replaceGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func replacePOST(w http.ResponseWriter, r *http.Request) {

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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}

}

func replaceGET(w http.ResponseWriter, r *http.Request) {

	//получим ID из запроса
	idLong := r.URL.Path[1:]
	id := strings.Split(idLong, "/")[0]
	if id == " " {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//проверим по ID ссылку
	if longURL, ok := short2long[id]; ok {
		w.Header().Set("Location", longURL)
		// } else {
		// 	longURL = "https://practicum.yandex.ru/"
		// 	w.Header().Set("Location", longURL)
		// 	short2long[id] = longURL
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(longURL))
		// http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
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
