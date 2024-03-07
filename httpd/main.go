package httpd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"plants/config"
	"plants/plants"
	"plants/store"
	"sync"
	"time"
)

func NewServer(logger *slog.Logger, plantStore store.Store) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger, plantStore)
	// NOTE: handler doesnt care about what mux we use, so we define it as interface
	var handler http.Handler = mux

	reqLoggerMiddleware := newLoggerMiddleware(logger)
	// NOTE: youd add more middleware in the chain here
	handler = reqLoggerMiddleware(handler)

	return handler
}

func Run(ctx context.Context) error {
	logger := slog.Default()
	// NOTE: realistically this wouldnt be an in-memory array,
	// but a DB implementation of store.Store interface
	s := store.NewMemoryStore([]plants.Plant{})
	srv := NewServer(logger, s)

	cfg := config.NewDefaultServer()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, string(cfg.Port)),
		Handler: srv,
	}

	go func() {
		logger.Info(fmt.Sprintf("listening to requests on %s\n", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx.Done()
		// TODO: figure out if this is broken
		shutdownCtx := context.Background()
		// allow 10 seconds to shut down gracefully
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}

func addRoutes(mux *http.ServeMux, logger *slog.Logger, plantStore store.Store) {
	root := "/api/v1"

	// NOTE: you can add specific middleware to each route here
	adminOnly := newAdminOnlyMiddleware("supersecret")
	mux.Handle(http.MethodGet+" "+root+"/plants/", handleListPlants(logger, plantStore))
	mux.Handle(http.MethodPost+" "+root+"/plants/", adminOnly(handleCreatePlant(logger, plantStore)))
	mux.Handle(http.MethodGet+" "+root+"/plants/{id}/", handleGetPlant(logger, plantStore))
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	return v, nil
}

type Validator interface {
	Valid() (problems map[string]string)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := decode[T](r)
	if err != nil {
		return v, nil, err
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}
