package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/BelyaevEI/shortener/internal/config"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/storage"
	"github.com/BelyaevEI/shortener/internal/utils"
)

func ReplacePOST() http.HandlerFunc {

	post := func(w http.ResponseWriter, r *http.Request) {

		var (
			LongShortUrl models.StorageURL
			short_id     string
		)
		// Открываем файл на чтение/запись
		f := storage.NewStorage()

		defer f.Close()

		storage := f.ReadAllURLS()

		//Считать из тела запроса строку URL
		longURL, err := io.ReadAll(r.Body)
		if err != nil || string(longURL) == " " {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if short_id = utils.TryFoundShortUrl(longURL, storage); short_id != " " {

			short_url := config.ShortURL + "/" + short_id
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(short_url))

		} else {

			short_id = utils.GenerateRandomString(8)
			LongShortUrl.OriginalURL = string(longURL)
			LongShortUrl.ShortURL = short_id
			f.WriteURL(&LongShortUrl) // Запись новой пары в файл

			short_url := config.ShortURL + "/" + short_id

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(short_url))
		}

	}
	return http.HandlerFunc(post)
}

func PostAPI() http.HandlerFunc {
	post := func(w http.ResponseWriter, r *http.Request) {
		var (
			req          models.Request
			shortURL     string
			buf          bytes.Buffer
			LongShortUrl models.StorageURL
			short_id     string
		)

		// Открываем файл на чтение/запись
		f := storage.NewStorage()

		defer f.Close()

		storage := f.ReadAllURLS()

		// читаем тело запроса
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON
		if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		longURL := req.URL

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if short_id = utils.TryFoundShortUrl([]byte(longURL), storage); short_id == " " {

			short_id = utils.GenerateRandomString(8)
			LongShortUrl.OriginalURL = longURL
			LongShortUrl.ShortURL = short_id
			f.WriteURL(&LongShortUrl) // Запись новой пары в файл
		}

		shortURL = config.ShortURL + "/" + short_id

		// заполняем модель ответа
		resp := models.Response{
			Result: shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		//сериализуем ответ сервера
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
	return http.HandlerFunc(post)
}

func ReplaceGET() http.HandlerFunc {

	get := func(w http.ResponseWriter, r *http.Request) {
		var id string

		// Открываем файл на чтение/запись
		f := storage.NewStorage()

		defer f.Close()

		storage := f.ReadAllURLS()

		//получим ID из запроса
		short_id := r.URL.Path[1:]
		// idLong := r.URL.Query().Get("id") не работает почему-то, не забудь разобраться

		if strings.ContainsRune(short_id, '/') {
			id = strings.Split(short_id, "/")[0]
			if id == " " {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			id = short_id
			if id == " " {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// Проверим, есть ли в файле нужная ссылка
		// если ее нет, отправляем 400 пользователю
		if origin_url := utils.TryFoundOrigUrl(id, storage); origin_url != " " {
			w.Header().Set("Location", origin_url)
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(origin_url))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	return http.HandlerFunc(get)
}
