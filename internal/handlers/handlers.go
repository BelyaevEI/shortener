package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	cookies "github.com/BelyaevEI/shortener/internal/cookie"
	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/storages/storage"
	"github.com/BelyaevEI/shortener/internal/utils"
)

type Handlers struct {
	shortURL string
	storage  *storage.Storage
	logger   *logger.Logger
}

func New(shortURL string, storage *storage.Storage, log *logger.Logger) Handlers {
	return Handlers{
		shortURL: shortURL,
		storage:  storage,
		logger:   log,
	}
}

func (h *Handlers) ReplacePOST(w http.ResponseWriter, r *http.Request) {

	var (
		shortid  string
		shortURL string
		status   int
	)

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cookie, _ := r.Cookie("Token")

	userID, _ := cookies.GetUserId(cookie.Value)

	//Считаем из тела запроса строку URL
	longURL, err := io.ReadAll(r.Body)
	if err != nil || string(longURL) == " " {
		h.logger.Log.Error("Empty body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование ссылки
	shortid, err = h.storage.GetShortenURL(ctx, string(longURL), userID)
	if err != nil {
		h.logger.Log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(shortid) == 0 {
		shortid = utils.GenerateRandomString(8)
		status = http.StatusCreated
		err = h.storage.SaveURL(ctx, shortid, string(longURL), userID)
		if err != nil {
			h.logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		status = http.StatusConflict
	}

	shortURL = h.shortURL + "/" + shortid
	utils.Response(w, "Content-Type", "text/plain", shortURL, status)
}

func (h *Handlers) PostAPI(w http.ResponseWriter, r *http.Request) {

	var (
		req      models.Request
		shortURL string
		buf      bytes.Buffer
		shortid  string
		status   int
	)

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cookie, _ := r.Cookie("Token")

	userID, _ := cookies.GetUserId(cookie.Value)

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.logger.Log.Error(err)
		return
	}

	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		h.logger.Log.Error(err)
		return
	}

	longURL := req.URL

	// Проверяем существование ссылки
	shortid, err = h.storage.GetShortenURL(ctx, longURL, userID)
	if err != nil {
		h.logger.Log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(shortid) == 0 {
		shortid = utils.GenerateRandomString(8)
		status = http.StatusCreated
		err = h.storage.SaveURL(ctx, shortid, longURL, userID)
		if err != nil {
			h.logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		status = http.StatusConflict
	}

	shortURL = h.shortURL + "/" + shortid

	// заполняем модель ответа
	resp := models.Response{
		Result: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		h.logger.Log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handlers) ReplaceGET(w http.ResponseWriter, r *http.Request) {

	var id string

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cookie, _ := r.Cookie("Token")

	userID, _ := cookies.GetUserId(cookie.Value)

	//получим ID из запроса
	shortid := r.URL.Path[1:]

	if strings.ContainsRune(shortid, '/') {
		id = strings.Split(shortid, "/")[0]
		if len(id) == 0 {
			h.logger.Log.Info("Empty id in Get request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		id = shortid
		if len(id) == 0 {
			h.logger.Log.Info("Empty id in Get request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Проверяем существование ссылки
	originURL, err := h.storage.GetOriginalURL(ctx, id, userID)
	if err != nil || len(originURL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Log.Infoln(err, originURL)
		return
	}

	utils.Response(w, "Location", originURL, originURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) PingDB(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := h.storage.Ping(ctx); err != nil {
		h.logger.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) PostAPIBatch(w http.ResponseWriter, r *http.Request) {
	var (
		batchinput  []models.Batch
		batchoutput []models.Batch
		shortid     string
		shortURL    string
	)

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cookie, _ := r.Cookie("Token")

	userID, _ := cookies.GetUserId(cookie.Value)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Log.Error("Error read body request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &batchinput)
	if err != nil {
		h.logger.Log.Error("Error deserialization", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range batchinput {

		shortid, err = h.storage.GetShortenURL(ctx, v.OriginalURL, userID)
		if err != nil {
			h.logger.Log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(shortid) == 0 {
			shortid = utils.GenerateRandomString(8)
			err = h.storage.SaveURL(ctx, shortid, string(v.OriginalURL), userID)
			if err != nil {
				h.logger.Log.Error("Error save data", err)
			}
		}

		shortURL = h.shortURL + "/" + shortid

		// заполняем модель ответа
		resp := models.Batch{
			CorrelationID: v.CorrelationID,
			ShortURL:      shortURL,
		}

		batchoutput = append(batchoutput, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(batchoutput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Log.Error("Error serialization", err)
		return
	}
}

func (h *Handlers) GetAllUrlsUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cookie, _ := r.Cookie("Token")

	userID, _ := cookies.GetUserId(cookie.Value)

	// находим все ссылки, которые сокращал данный пользователь
	// если таковых нет, ответ короткий
	allURLS, err := h.storage.GetUrlsUser(ctx, userID)
	if err != nil || len(allURLS) == 0 {
		h.logger.Log.Error(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//сериализуем ответ сервера
	enc := json.NewEncoder(w)
	if err := enc.Encode(allURLS); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Log.Error("Error serialization", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
