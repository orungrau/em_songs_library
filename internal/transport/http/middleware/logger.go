package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	log zerolog.Logger
}

func NewLoggerMiddleware(log zerolog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{log: log.With().Str("module", "http-transport").Logger()}
}

func (m *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		// TODO: Get span-id from headers
		ctx = context.WithValue(ctx, "span-id", uuid.New().String())
		ctx = context.WithValue(ctx, "duration", start)
		r = r.WithContext(ctx)

		logger := m.log.With().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Logger()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			if rec := recover(); rec != nil {
				logger.Error().
					Interface("panic", rec).
					Int("status", http.StatusInternalServerError).
					Msg("Recovered from panic")
				http.Error(lrw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if lrw.statusCode >= 400 {
				logger.Error().
					Ctx(ctx).
					Int("status", lrw.statusCode).
					Msg("HTTP Request completed with error")
			} else {
				logger.Info().
					Ctx(ctx).
					Int("status", lrw.statusCode).
					Msg("HTTP Request completed successfully")
			}
		}()

		next.ServeHTTP(lrw, r)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}
