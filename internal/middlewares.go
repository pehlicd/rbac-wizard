package internal

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func (l App) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{w: w, status: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		level := zerolog.InfoLevel
		if ww.status >= 500 {
			level = zerolog.ErrorLevel
		} else if ww.status >= 400 {
			level = zerolog.WarnLevel
		}

		l.Logger.WithLevel(level).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.status).
			Str("remote_addr", r.RemoteAddr).
			Dur("duration", duration).
			Msg("")
	})
}

type responseWriter struct {
	w           http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.status = statusCode
		rw.w.WriteHeader(statusCode)
		rw.wroteHeader = true
	}
}
