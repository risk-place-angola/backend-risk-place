package middleware

import (
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		userID, _ := r.Context().Value(util.UserIDCtxKey).(string)

		slog.Info("Request processed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.Int("status", rw.status),
			slog.String("duration", duration.String()),
			slog.String("remote_ip", r.RemoteAddr),
			slog.String("user_id", userID),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	if lw.wroteHeader {
		return
	}
	lw.status = code
	lw.wroteHeader = true
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	if !lw.wroteHeader {
		lw.WriteHeader(http.StatusOK)
	}
	return lw.ResponseWriter.Write(b)
}
