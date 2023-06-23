package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	t.Run(test1.name, func(t *testing.T) {
		//Создаем тело запроса
		requestBody := strings.NewReader("https://practicum.yandex.ru/ ")

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "text/plain")

		//Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		//Обрабатываем запрос
		r := ReplacePOST()
		r(responseRecorder, request)

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
		requestBody := strings.NewReader(" ")

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodPost, "/", requestBody)

		//Устанавливаем заголовок
		request.Header.Set("Content-Type", "text/plain")

		//Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		//Обрабатываем запрос
		r := ReplacePOST()
		r(responseRecorder, request)

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
	type test struct {
		name string
		code int
	}

	//Сформируем варианты для тестирования
	test1 := test{name: "Empty ID in URL", code: http.StatusBadRequest}

	t.Run(test1.name, func(t *testing.T) {

		//Создаем сам запрос
		request := httptest.NewRequest(http.MethodGet, "/asd/", nil)

		//Создаем рекордер для записи ответа
		responseRecorder := httptest.NewRecorder()

		//Обрабатываем запрос
		r := ReplaceGET()

		// ReplaceGET(responseRecorder, request)
		r(responseRecorder, request)

		//Получаем ответ
		result := responseRecorder.Result()

		defer result.Body.Close()

		//Делаем проверки
		//Проверка ответа сервера
		assert.Equal(t, test1.code, result.StatusCode)

	})

}
