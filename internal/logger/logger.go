package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logGet := func(w http.ResponseWriter, r *http.Request) {

		// создаём предустановленный регистратор zap
		logger, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic(err)
		}

		defer logger.Sync()

		// делаем регистратор SugaredLogger
		sugar = *logger.Sugar()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		//Время запуска
		start := time.Now()

		// эндпоинт
		uri := r.RequestURI

		// метод запроса
		method := r.Method

		// обслуживание оригинального запроса
		// внедряем реализацию http.ResponseWriter
		h.ServeHTTP(&lw, r)

		//время выполнения
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logGet)
}
