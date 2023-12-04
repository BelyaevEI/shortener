package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "net/http/pprof"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/storages/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplacePOST(t *testing.T) {

	//Структура запроса
	type want struct {
		code        int
		contentType string
	}

	type test struct {
		name string
		want want
	}

	//Сформируем варианты для тестирования
	test1 := test{name: "Simple test POST request",
		want: want{
			code:        http.StatusCreated,
			contentType: "text/plain",
		},
	}

	test2 := test{name: "Empty link in body",
		want: want{
			code: http.StatusBadRequest,
		},
	}

	//Создаем логгер
	log := logger.New()

	storage := storage.Init(" ", "", log)

	//Создаем обьект handle
	h := New("http://localhost:8080", "", storage, log)

	t.Run(test1.name, func(t *testing.T) {

		//Создаем тело запроса
		requestBody := strings.NewReader("https://practicum.yandex.ru/")

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "text/plain")

		//Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		//Обрабатываем запрос
		h.ReplacePOST(responseRecorder, request)

		//Получаем ответ
		result := responseRecorder.Result()

		//Делаем проверки
		//Проверка ответа сервера
		assert.Equal(t, test1.want.code, result.StatusCode)

		//Проверка типа контента
		assert.Equal(t, result.Header.Get("Content-Type"), test1.want.contentType)

		//Получаем тело ответа
		resBody, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		//Проверка ответа без ошибок
		require.NoError(t, err)

		//Проверка тела ответа на пустоту
		assert.NotEmpty(t, string(resBody))

	})

	t.Run(test2.name, func(t *testing.T) {

		//Создаем тело запроса
		requestBody := strings.NewReader("")

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "text/plain")

		//Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		//Обрабатываем запрос
		h.ReplacePOST(responseRecorder, request)

		//Получаем ответ
		result := responseRecorder.Result()

		defer result.Body.Close()
		//Делаем проверки
		//Проверка ответа сервера
		assert.Equal(t, test2.want.code, result.StatusCode)

	})
}

func TestReplaceGET(t *testing.T) {
	//Структура запроса
	type want struct {
		code        int
		contentType string
		url         string
	}

	type test struct {
		name string
		want want
	}

	// Сформируем варианты для тестирования
	dataTest := test{name: "Simple test GET request",
		want: want{
			code:        http.StatusTemporaryRedirect,
			contentType: "text/plain",
			url:         "https://practicum.yandex.ru/",
		},
	}

	ctx := context.Background()

	// Создаем логгер
	log := logger.New()

	storage := storage.Init(" ", "", log)
	storage.SaveURL(ctx, "TESTURL", "https://practicum.yandex.ru/", 1234)

	// Создаем обьект handle
	h := New("http://localhost:8080", "", storage, log)

	t.Run(dataTest.name, func(t *testing.T) {

		//Создаем тело запроса
		requestBody := strings.NewReader("")

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodGet, "/TESTURL", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "text/plain")

		// Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		// Обрабатываем запрос
		h.ReplaceGET(responseRecorder, request)

		// Получаем ответ
		result := responseRecorder.Result()

		assert.Equal(t, dataTest.want.code, result.StatusCode)

		//Получаем тело ответа
		resBody, err := io.ReadAll(result.Body)
		defer result.Body.Close()

		//Проверка ответа без ошибок
		require.NoError(t, err)

		//Проверка тела ответа
		assert.Equal(t, dataTest.want.url, string(resBody))
	})
}

func TestPostAPI(t *testing.T) {
	// Структура запроса
	type want struct {
		code        int
		contentType string
	}

	type res struct {
		Result string `json:"result"`
	}

	type test struct {
		name string
		want want
	}

	// Сформируем варианты для тестирования
	dataTest := test{
		name: "Post API test",
		want: want{code: http.StatusCreated,
			contentType: "application/json"},
	}

	// Создаем логгер
	log := logger.New()

	storage := storage.Init(" ", "", log)

	// Создаем обьект handle
	h := New("http://localhost:8080", "", storage, log)

	t.Run(dataTest.name, func(t *testing.T) {

		var res res

		r := models.Request{URL: "https://practicum.yandex.ru"}

		req, _ := json.Marshal(r)

		//Создаем тело запроса
		requestBody := strings.NewReader(string(req))

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "json/application")

		// Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		// Обрабатываем запрос
		h.PostAPI(responseRecorder, request)

		// Получаем ответ
		result := responseRecorder.Result()

		//Получаем тело ответа
		resBody, _ := io.ReadAll(result.Body)
		defer result.Body.Close()

		json.Unmarshal(resBody, &res)

		//Проверка тела ответа на пустоту
		assert.NotEmpty(t, res)

		//Проверка ответа сервера
		assert.Equal(t, dataTest.want.code, result.StatusCode)

	})
}

func TestPostAPIBatch(t *testing.T) {
	// Структура запроса
	type want struct {
		code        int
		contentType string
	}

	type test struct {
		name string
		want want
	}

	// Сформируем варианты для тестирования
	dataTest := test{
		name: "Post API batch test",
		want: want{code: http.StatusCreated,
			contentType: "application/json"},
	}

	ctx := context.Background()

	// Создаем логгер
	log := logger.New()

	storage := storage.Init(" ", "", log)
	storage.SaveURL(ctx, "TESTURL", "https://practicum.yandex.ru/", 1234)
	storage.SaveURL(ctx, "TESTURL2", "https://yandex.ru/", 1234)

	// Создаем обьект handle
	h := New("http://localhost:8080", "", storage, log)

	t.Run(dataTest.name, func(t *testing.T) {
		var res models.Batch
		r := make([]models.Batch, 0)

		r = append(r, models.Batch{CorrelationID: "1",
			OriginalURL: "https://practicum.yandex.ru"})

		r = append(r, models.Batch{CorrelationID: "2",
			OriginalURL: "https://yandex.ru"})

		req, _ := json.Marshal(r)

		//Создаем тело запроса
		requestBody := strings.NewReader(string(req))

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "json/application")

		// Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		// Обрабатываем запрос
		h.PostAPIBatch(responseRecorder, request)

		// Получаем ответ
		result := responseRecorder.Result()

		//Получаем тело ответа
		resBody, _ := io.ReadAll(result.Body)
		defer result.Body.Close()

		err := json.Unmarshal(resBody, &res)
		if err != nil {
			fmt.Println("Empty ")
		}
		// //Проверка тела ответа на пустоту
		// assert.NotEmpty(t, res)

		//Проверка ответа сервера
		assert.Equal(t, dataTest.want.code, result.StatusCode)
	})
}

func BenchmarkReplacePOST(b *testing.B) {

	//Создаем логгер
	log := logger.New()

	storage := storage.Init(" ", "", log)

	//Создаем обьект handle
	h := New("http://localhost:8080", "", storage, log)

	//Создаем тело запроса
	requestBody := strings.NewReader("https://practicum.yandex.ru/")

	//Создаем сам запрос
	request := httptest.NewRequest(http.MethodPost, "/", requestBody)

	//Устанавливаем заголовок
	request.Header.Set("Content-Type", "text/plain")

	//Создаем рекордер для записи ответа
	responseRecorder := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {

		//Обрабатываем запрос
		h.ReplacePOST(responseRecorder, request)
	}
}
