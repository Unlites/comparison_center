package middleware

import (
	"net/http"

	"time"

	"github.com/Unlites/comparison_center/backend/internal/adapters/handlers/http/v1/metrics"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedResponseWriter{w, http.StatusOK}
		next.ServeHTTP(ww, r)
		metrics.ObserveRequest(time.Since(start), ww.statusCode)
	})
}
