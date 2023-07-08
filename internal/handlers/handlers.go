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

var short2long = make(map[string]string) //Словарь для получения полного URL по короткому
var long2short = make(map[string]string) //Словарь для получения короткого URL по полному

type Handlers struct {
	FileStoragePath string
	ShortURL        string
	s               storage.Storage
}

func New(cfg config.Parameters, s *storage.Storage) Handlers {
	return Handlers{
		FileStoragePath: cfg.FileStoragePath,
		ShortURL:        cfg.ShortURL,
	}
}

func (h *Handlers) ReplacePOST(w http.ResponseWriter, r *http.Request) {

	var (
		LongShortURL models.StorageURL
		shortid      string
		shortURL     string
	)

	//Считаем из тела запроса строку URL
	longURL, err := io.ReadAll(r.Body)
	if err != nil || string(longURL) == " " {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(h.FileStoragePath) != 0 {

		//Читаем весь файл
		storage := h.s.ReadAllURLS()

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if shortid = utils.TryFoundShortURL(longURL, storage); shortid != "" {

			shortURL = h.ShortURL + "/" + shortid
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))

		} else {

			shortid = utils.GenerateRandomString(8)
			LongShortURL.OriginalURL = string(longURL)
			LongShortURL.ShortURL = shortid
			h.s.WriteURL(&LongShortURL) // Запись новой пары в файл

			shortURL = h.ShortURL + "/" + shortid

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
		h.s.Close()
	} else {
		if shortURL, ok := long2short[string(longURL)]; ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		} else {
			short := utils.GenerateRandomString(8)
			shortURL = h.ShortURL + "/" + short

			long2short[string(longURL)] = shortURL
			short2long[short] = string(longURL)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
}

func (h *Handlers) PostAPI(w http.ResponseWriter, r *http.Request) {
	var (
		req          models.Request
		shortURL     string
		buf          bytes.Buffer
		LongShortURL models.StorageURL
		shortid      string
	)

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

	if len(h.FileStoragePath) != 0 {

		//Читаем весь файл
		storage := h.s.ReadAllURLS()

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if shortid = utils.TryFoundShortURL([]byte(longURL), storage); shortid == "" {

			shortid = utils.GenerateRandomString(8)
			LongShortURL.OriginalURL = longURL
			LongShortURL.ShortURL = shortid
			h.s.WriteURL(&LongShortURL) // Запись новой пары в файл
		}

		shortURL = h.ShortURL + "/" + shortid

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

		h.s.Close()
	} else {
		if shortURL, ok := long2short[string(longURL)]; ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		} else {
			short := utils.GenerateRandomString(8)

			shortURL = h.ShortURL + "/" + short

			long2short[string(longURL)] = shortURL
			short2long[short] = string(longURL)

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(shortURL))
		}
	}
}

func (h *Handlers) ReplaceGET(w http.ResponseWriter, r *http.Request) {

	var id string

	//получим ID из запроса
	shortid := r.URL.Path[1:]

	// idLong := r.URL.Query().Get("id") не работает почему-то, не забудь разобраться
	if strings.ContainsRune(shortid, '/') {
		id = strings.Split(shortid, "/")[0]
		if len(id) == 0 { // Если заменить на len(id) == 0?
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		id = shortid
		if len(id) == 0 { // Если заменить на len(id) == 0?
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if len(h.FileStoragePath) != 0 {

		//Читаем весь файл
		storage := h.s.ReadAllURLS()

		// Проверим, есть ли в файле нужная ссылка
		// если ее нет, отправляем 400 пользователю
		if originURL := utils.TryFoundOrigURL(id, storage); originURL != "" {
			w.Header().Set("Location", originURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(originURL))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		h.s.Close()
	} else {
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
}
