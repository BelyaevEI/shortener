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
		ResponseWriter http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		RespData       *ResponseDatas
	}
)

// Header implements http.ResponseWriter.
func (LoggResponse) Header() http.Header {
	panic("unimplemented")
}

func (r LoggResponse) Write(b []byte) (int, error) {

	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.RespData.Size += size // захватываем размер
	return size, err
}

func (r LoggResponse) WriteHeader(statusCode int) {

	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.RespData.Status = statusCode // захватываем код статуса
}
