package http

import (
	"net/http"
	"time"

	hu "github.com/Unlites/comparison_center/backend/internal/utils/http"
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
		hu.ObserveRequest(time.Since(start), ww.statusCode)
	})
}
