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

func NewServer(logger *slog.Logger, config config.Server, plantStore store.Store) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, config, logger, plantStore)
	// NOTE: handler doesnt care about what mux we use, so we define it as interface
	var handler http.Handler = mux

	reqLoggerMiddleware := newLoggerMiddleware(logger)
	// NOTE: youd add more middleware in the chain here
	handler = reqLoggerMiddleware(handler)

	return handler
}

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	logger := slog.New(slog.NewJSONHandler(stdout, nil))
	slog.SetDefault(logger)

	cfg := config.NewDefaultServer()
	// you could get some env vars here using `getenv`
	// cfg.HOST = getenv("API_HOST")
	// or parse args

	// NOTE: realistically this wouldnt be an in-memory array,
	// but a DB implementation of store.Store interface
	s := store.NewMemoryStore([]plants.Plant{})

	httpHandler := NewServer(logger, cfg, s)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: httpHandler,
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
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("error shutting down: %s", err))
	}

	return nil
}

func addRoutes(mux *http.ServeMux, config config.Server, logger *slog.Logger, plantStore store.Store) {
	root := config.RootPrefix

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
