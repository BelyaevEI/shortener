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
			shortid      string
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
		if shortid = utils.TryFoundShortURL(longURL, storage); shortid != " " {

			short_url := config.ShortURL + "/" + shortid
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(short_url))

		} else {

			shortid = utils.GenerateRandomString(8)
			LongShortUrl.OriginalURL = string(longURL)
			LongShortUrl.ShortURL = shortid
			f.WriteURL(&LongShortUrl) // Запись новой пары в файл

			short_url := config.ShortURL + "/" + shortid

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
			shortid      string
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
		if shortid = utils.TryFoundShortURL([]byte(longURL), storage); shortid == " " {

			shortid = utils.GenerateRandomString(8)
			LongShortUrl.OriginalURL = longURL
			LongShortUrl.ShortURL = shortid
			f.WriteURL(&LongShortUrl) // Запись новой пары в файл
		}

		shortURL = config.ShortURL + "/" + shortid

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
		shortid := r.URL.Path[1:]
		// idLong := r.URL.Query().Get("id") не работает почему-то, не забудь разобраться

		if strings.ContainsRune(shortid, '/') {
			id = strings.Split(shortid, "/")[0]
			if id == " " {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			id = shortid
			if id == " " {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// Проверим, есть ли в файле нужная ссылка
		// если ее нет, отправляем 400 пользователю
		if originurl := utils.TryFoundOrigURL(id, storage); originurl != " " {
			w.Header().Set("Location", originurl)
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(originurl))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	return http.HandlerFunc(get)
}
