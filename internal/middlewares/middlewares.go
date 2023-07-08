package midlleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/BelyaevEI/shortener/internal/compres"
	"github.com/BelyaevEI/shortener/internal/logger"
	"go.uber.org/zap"
)

// Middleware - мидлварь сжатия
func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {

			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := compres.NewWriter(w)

			// меняем оригинальный http.ResponseWriter на новый
			ow = cw

			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := compres.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}
		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}

// Middleware - мидлварь логер
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// создаём предустановленный регистратор zap
		logg, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic(err)
		}

		defer logg.Sync()

		// делаем регистратор SugaredLogger
		sugar := *logg.Sugar()

		responseData := &logger.ResponseDatas{
			Status: 0,
			Size:   0,
		}

		lw := logger.LoggResponse{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			RespData:       responseData,
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
			"status", responseData.Status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.Size, // получаем перехваченный размер ответа
		)
	})
}
