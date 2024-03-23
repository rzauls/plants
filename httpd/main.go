package httpd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"plants/config"
	"plants/plants"
	"plants/store"
)

func NewApiHandler(logger *slog.Logger, config config.Server, plantStore store.Store) http.Handler {
	mux := http.NewServeMux()

	// NOTE: you can add specific middleware to each route here
	adminOnly := newAdminOnly("supersecret")

	mux.Handle("GET /health", handleHealth())
	mux.Handle("GET /plants/", handleListPlants(plantStore))
	mux.Handle("POST /plants/", adminOnly(handleCreatePlant(plantStore)))
	mux.Handle("GET /plants/{id}/", handleGetPlant(plantStore))

	root := http.NewServeMux()
	root.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	stack := newMiddlewareStack(
		newTracing(logger),
		newLogger(logger),
	)
	var handler http.Handler = root

	return stack(handler)
}

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	// NOTE: you would configure the slog.Logger appropriately here, this is just for running locally,
	// in a human-readable(-ish) form
	logger := slog.New(slog.NewTextHandler(stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg := config.FromEnv(getenv)

	// NOTE: realistically this wouldnt be an in-memory array,
	// but a DB implementation of store.Store interface
	s := store.NewMemoryStore([]plants.Plant{})

	handler := NewApiHandler(logger, cfg, s)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: handler,
	}

	go func() {
		logger.Info(fmt.Sprintf("listening to requests on %s", string(httpServer.Addr)))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("error listening: %s", err))
			return
		}
	}()

	<-ctx.Done()
	logger.Info("graceful shutdown")
	if err := httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		logger.Error(fmt.Sprintf("error shutting down: %s", err))
	}

	return nil
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
		return v, problems, fmt.Errorf("invalid input with %d error(-s)", len(problems))
	}

	return v, nil, nil
}
