package midllewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/BelyaevEI/shortener/internal/compres"
	cookies "github.com/BelyaevEI/shortener/internal/cookie"
	"github.com/BelyaevEI/shortener/internal/models"
	"github.com/BelyaevEI/shortener/internal/utils"
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

// func Gzip(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Header.Get(`Content-Encoding`) == `gzip` {
// 			gz, err := gzip.NewReader(r.Body)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}

// 			r.Body = gz
// 			defer gz.Close()
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// Middleware - работа с куки
func Cookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const keyID models.KeyID = "userID"

		cookie, err := r.Cookie("Token")

		if err != nil {

			if !errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// если кук нет, то сгенерируем
			userID := utils.GenerateUniqueID()
			cookies.NewCookie(w, userID)
			ctx := context.WithValue(r.Context(), keyID, userID)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		// если кук нет или валидация не прошла
		// генерируем новые куки по заданию
		if cookie == nil || !cookies.Validation(cookie.Value) {
			userID := utils.GenerateUniqueID()
			cookies.NewCookie(w, userID)
			ctx := context.WithValue(r.Context(), keyID, userID)
			h.ServeHTTP(w, r.WithContext(ctx))
		}

		h.ServeHTTP(w, r)

	})
}
