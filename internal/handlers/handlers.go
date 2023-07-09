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

// var short2long = make(map[string]string) //Словарь для получения полного URL по короткому
// var long2short = make(map[string]string) //Словарь для получения короткого URL по полному

type Handlers struct {
	FileStoragePath string
	ShortURL        string
	Config          config.Parameters
	short2long      map[string]string
	long2short      map[string]string
}

func New(cfg config.Parameters) Handlers {
	return Handlers{
		FileStoragePath: cfg.FileStoragePath,
		ShortURL:        cfg.ShortURL,
		Config:          cfg,
		short2long:      make(map[string]string),
		long2short:      make(map[string]string),
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

	if h.FileStoragePath != " " {

		//Работа с файлом
		s := storage.New(h.Config)

		//Читаем весь файл
		storage := s.ReadAllURLS()

		defer s.Close()

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if shortid = utils.TryFoundShortURL(longURL, storage); shortid != "" {

			shortURL = h.ShortURL + "/" + shortid
			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)

		} else {

			shortid = utils.GenerateRandomString(8)

			LongShortURL.OriginalURL, LongShortURL.ShortURL = string(longURL), shortid
			s.WriteURL(&LongShortURL) // Запись новой пары в файл

			shortURL = h.ShortURL + "/" + shortid

			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
		}

	} else {

		if shortURL, ok := h.long2short[string(longURL)]; ok {
			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
		} else {
			short := utils.GenerateRandomString(8)
			shortURL = h.ShortURL + "/" + short

			h.long2short[string(longURL)] = shortURL
			h.short2long[short] = string(longURL)

			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
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

	if h.FileStoragePath != " " {

		//Работа с файлом
		s := storage.New(h.Config)

		//Читаем весь файл
		storage := s.ReadAllURLS()

		defer s.Close()

		// Проверяем есть ли в файле ссылка, если нет, то сгенерируем,
		// запишем в файл и отправим пользвоателю
		if shortid = utils.TryFoundShortURL([]byte(longURL), storage); shortid == "" {

			shortid = utils.GenerateRandomString(8)
			LongShortURL.OriginalURL, LongShortURL.ShortURL = longURL, shortid
			s.WriteURL(&LongShortURL) // Запись новой пары в файл
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

	} else {
		if shortURL, ok := h.long2short[string(longURL)]; ok {
			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
		} else {
			short := utils.GenerateRandomString(8)

			shortURL = h.ShortURL + "/" + short

			h.long2short[string(longURL)] = shortURL
			h.short2long[short] = string(longURL)

			utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
		}
	}
}

func (h *Handlers) ReplaceGET(w http.ResponseWriter, r *http.Request) {

	var id string

	//получим ID из запроса
	shortid := r.URL.Path[1:]

	if strings.ContainsRune(shortid, '/') {
		id = strings.Split(shortid, "/")[0]
		if len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		id = shortid
		if len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if h.FileStoragePath != " " {

		//Работа с файлом
		s := storage.New(h.Config)

		//Читаем весь файл
		storage := s.ReadAllURLS()

		defer s.Close()

		// Проверим, есть ли в файле нужная ссылка
		// если ее нет, отправляем 400 пользователю
		if originURL := utils.TryFoundOrigURL(id, storage); originURL != "" {
			utils.Response(w, "Location", originURL, originURL, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	} else {
		//проверим по ID ссылку
		if longURL, ok := h.short2long[id]; ok {
			utils.Response(w, "Location", longURL, longURL, http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
