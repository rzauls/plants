package httpd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"plants/log"
	"time"

	"github.com/rs/xid"
)

type Middleware func(http.Handler) http.Handler

func newMiddlewareStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}

		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func newLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	slog.SetDefault(logger)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			requestLogger := log.LoggerFromCtx(r.Context())
			requestLogger.Info(
				fmt.Sprintf("%s %s", r.Method, r.URL.String()),
				slog.Int("statusCode", wrapped.statusCode),
				slog.String("duration", time.Since(start).String()),
			)
		})
	}
}

func newAdminOnly(_ string) func(next http.Handler) http.Handler {
	// NOTE: this doesnt actually do anything, its just an example of how a middleware would get its dependencies
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: some auth handler thing could go here (like a JWT check)
			next.ServeHTTP(w, r)
		})
	}
}

// NOTE: use typed strings for context keys so they cannot collide on accident,
// however in this case we are using constants as keys so that wouldnt be possible anyway
type traceCtxKey string

const CONTEXT_TRACE_ID traceCtxKey = "ctx.trace"

func newTracing(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// generate request ID
			requestID := xid.New().String()
			ctx := context.WithValue(r.Context(), CONTEXT_TRACE_ID, requestID)

			// add request ID to all child logs
			scopedLogger := logger.With(slog.String("traceId", requestID))
			ctx = context.WithValue(ctx, log.CONTEXT_LOGGER, scopedLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
