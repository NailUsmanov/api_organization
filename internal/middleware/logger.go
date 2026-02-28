// Package middleware предоставляет промежуточные обработчики (middleware) для HTTP-сервера.
package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger возвращает middleware для логирования всех входящих HTTP-запросов.
func Logger(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": time.Since(start),
			}).Info("request handled")
		})
	}
}
