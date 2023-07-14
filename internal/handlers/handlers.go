package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/storages/storage"
	"github.com/BelyaevEI/shortener/internal/utils"
)

type Handlers struct {
	shortURL string
	storage  *storage.Storage
}

func New(shortURL string, storage *storage.Storage) Handlers {
	return Handlers{
		shortURL: shortURL,
		storage:  storage,
	}
}

func (h *Handlers) ReplacePOST(w http.ResponseWriter, r *http.Request) {

	var (
		shortid  string
		shortURL string
	)

	//Считаем из тела запроса строку URL
	longURL, err := io.ReadAll(r.Body)
	if err != nil || string(longURL) == " " {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование ссылки
	if shortid = h.storage.GetURL(string(longURL)); shortid == "" {
		shortid = utils.GenerateRandomString(8)
		err := h.storage.SaveURL(shortid, string(longURL))
		if err != nil {
			log.Fatal(err)
		}
	}

	shortURL = h.shortURL + "/" + shortid
	utils.Response(w, "Content-Type", "text/plain", shortURL, http.StatusCreated)
}

func (h *Handlers) PostAPI(w http.ResponseWriter, r *http.Request) {

	var (
		req      models.Request
		shortURL string
		buf      bytes.Buffer
		shortid  string
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

	if shortid = h.storage.GetURL(longURL); shortid == "" {
		shortid = utils.GenerateRandomString(8)
		err := h.storage.SaveURL(shortid, longURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	shortURL = h.shortURL + "/" + shortid

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

	if originURL := h.storage.GetURL(id); originURL != "" {
		utils.Response(w, "Location", originURL, originURL, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handlers) PingDB(w http.ResponseWriter, r *http.Request) {

	if err := h.storage.Ping(); err != nil {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &batchinput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range batchinput {

		if shortid = h.storage.GetURL(v.OriginalURL); shortid == "" {
			shortid = utils.GenerateRandomString(8)
			err := h.storage.SaveURL(shortid, string(v.OriginalURL))
			if err != nil {
				log.Fatal(err)
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
		return
	}
}
