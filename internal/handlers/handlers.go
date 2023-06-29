package handlers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/models"
)

// Вынести в отдельный каталог?
var short2long = make(map[string]string) //Словарь для получения полного URL по короткому
var long2short = make(map[string]string) //Словарь для получения короткого URL по полному

func ReplacePOST() http.HandlerFunc {
	post := func(w http.ResponseWriter, r *http.Request) {
		//Считать из тела запроса строку URL
		longURL, err := io.ReadAll(r.Body)
		if err != nil || string(longURL) == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//Не проходит тесты первого инкремента
		// //проверим что ссылка передается в Header с нужным типом
		// contentType := r.Header.Get("Content-Type")
		// if contentType != "text/plain" {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		//Проверим наличие короткой ссылки по длинной, если ее нет
		//то сгенерируем и запишем в словарь
		if shortURL, ok := long2short[string(longURL)]; ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		} else {
			short := generateRandomString(8)
			// shortURL = "http://localhost:8080/" + short

			shortURL = config.ShortURL + "/" + short

			long2short[string(longURL)] = shortURL
			short2long[short] = string(longURL)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
	return http.HandlerFunc(post)
}

func PostAPI() http.HandlerFunc {
	post := func(w http.ResponseWriter, r *http.Request) {
		var (
			req      models.Request
			shortURL string
		)

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//проверим на пустоту приходящую ссылку
		if req.Url == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		longURL := req.Url

		//Проверим наличие короткой ссылки по длинной, если ее нет
		//то сгенерируем и запишем в словарь
		if s, ok := long2short[string(longURL)]; !ok {

			short := generateRandomString(8)

			shortURL = config.ShortURL + "/" + short

			long2short[string(longURL)] = shortURL
			short2long[short] = string(longURL)

		} else {
			shortURL = s
		}

		// заполняем модель ответа
		resp := models.Response{
			Result: shortURL,
		}

		//сериализуем ответ сервера
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

	}
	return http.HandlerFunc(post)
}

func ReplaceGET() http.HandlerFunc {

	get := func(w http.ResponseWriter, r *http.Request) {
		var id string

		//получим ID из запроса
		// idLong := r.URL.Path[1:]
		idLong := r.URL.Query().Get("id")

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
	return http.HandlerFunc(get)
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
