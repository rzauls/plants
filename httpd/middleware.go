package httpd

import (
	"fmt"
	"log/slog"
	"net/http"
)

func newLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf(
				"%s %s",
				r.Method,
				r.URL.String(),
			))
			next.ServeHTTP(w, r)
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
