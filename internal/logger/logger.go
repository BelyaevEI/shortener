package logger

import (
	"net/http"
)

// var sugar zap.SugaredLogger

type (
	// берём структуру для хранения сведений об ответе
	ResponseDatas struct {
		Status int
		Size   int
	}

	// добавляем реализацию http.ResponseWriter
	LoggResponse struct {
		Writer   http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		RespData *ResponseDatas
	}
)

func (r LoggResponse) Header() http.Header {
	return r.Writer.Header()
}

func (r LoggResponse) Write(b []byte) (int, error) {

	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.Writer.Write(b)
	r.RespData.Size += size // захватываем размер
	return size, err
}

func (r LoggResponse) WriteHeader(statusCode int) {

	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.Writer.WriteHeader(statusCode)
	r.RespData.Status = statusCode // захватываем код статуса
}
