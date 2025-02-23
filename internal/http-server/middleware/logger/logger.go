// Package logger предоставляет middleware для логирования HTTP-запросов.
package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// New создает middleware для логирования информации о HTTP-запросах.
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("middleware logger enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// Обертка для получения статуса ответа и размера данных
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Фиксируем время начала обработки запроса
			t1 := time.Now()
			defer func() {
				// Логируем продолжительность выполнения запроса
				log.With(slog.String("duration", time.Since(t1).String()))
			}()

			// Передаем управление следующему обработчику в цепочке
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
