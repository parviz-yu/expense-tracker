package handlers

import (
	"net/http"
	"time"

	"github.com/parviz-yu/expense-tracker/pkg/logger"

	"github.com/go-chi/chi/v5/middleware"
)

func NewMWLogger(log logger.LoggerI) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := logger.With(
			log,
			logger.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := logger.With(
				log,
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("remote_addr", r.Header.Get("X-Real-IP")),
				logger.String("user_agent", r.UserAgent()),
				logger.String("request_id", r.Header.Get("X-Request-Id")),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					logger.Int("status", ww.Status()),
					logger.Int("bytes", ww.BytesWritten()),
					logger.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
