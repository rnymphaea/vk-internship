package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"vk-internship/internal/logger"
)

func LoggingMiddleware(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			log.Infof("HTTP request",
				map[string]interface{}{
					"method":     r.Method,
					"path":       r.URL.Path,
					"remote_ip":  r.RemoteAddr,
					"status":     ww.Status(),
					"bytes":      ww.BytesWritten(),
					"duration":   duration.String(),
					"user_agent": r.UserAgent(),
					"request_id": middleware.GetReqID(r.Context()),
				},
			)
		})
	}
}
