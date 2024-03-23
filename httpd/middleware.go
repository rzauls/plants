package httpd

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func NewMiddlewareStack(middlewares ...Middleware) Middleware {
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

func newLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				fmt.Sprintf("%s %s", r.Method, r.URL.String()),
				slog.Int("statusCode", wrapped.statusCode),
			)
		})
	}
}

func newAdminOnlyMiddleware(authHandler string) func(next http.Handler) http.Handler {
	// NOTE: this doesnt actually do anything, its just an example of how a middleware would get its dependencies
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: some auth handler thing could go here (like a JWT check)
			next.ServeHTTP(w, r)
		})
	}
}

func newTracingMiddleware(tracingHandler int) func(next http.Handler) http.Handler {
	// NOTE: another middleware that doesnt actually do anything, but is used as an example for stacking/chaining
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
