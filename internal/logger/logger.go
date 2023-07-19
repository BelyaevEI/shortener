package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

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

	Logger struct {
		Log zap.SugaredLogger
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

func New() *Logger {
	// создаём предустановленный регистратор zap
	logg, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}

	defer logg.Sync()

	// делаем регистратор SugaredLogger
	sugar := *logg.Sugar()

	return &Logger{Log: sugar}
}

func (l *Logger) Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseData := &ResponseDatas{
			Status: 0,
			Size:   0,
		}

		lw := LoggResponse{
			Writer:   w, // встраиваем оригинальный http.ResponseWriter
			RespData: responseData,
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
		l.Log.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.Status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.Size, // получаем перехваченный размер ответа
		)

	})
}
