package httpd

import (
	"errors"
	"fmt"
	"net/http"
	"plants/log"
	"plants/plants"
	"plants/store"
)

// TODO: This `encode` approach doesnt rly work well with error reporting
// there probabbly is a nicer way to report json marshalling errors

func handleHealth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = encode(w, r, http.StatusOK, "")
	})
}

func handleListPlants(plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.LoggerFromCtx(ctx)
		plts, err := plantStore.List(ctx)
		if err != nil {
			err = fmt.Errorf("retrieve all plants: %w", err)
			logger.Error(err.Error())
			_ = encode(w, r, http.StatusInternalServerError, newHttpError(err))
			return
		}

		if len(plts) == 0 {
			plts = make([]plants.Plant, 0)
		}

		_ = encode(w, r, http.StatusOK, plts)
	})
}

func handleGetPlant(plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.LoggerFromCtx(ctx)
		id := r.PathValue("id")
		if id == "" {
			err := errors.New("id is required in path parameters")
			logger.Error(err.Error())
			_ = encode(w, r, http.StatusUnprocessableEntity, newHttpError(err))
			return
		}

		plant, err := plantStore.Find(ctx, id)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.As(err, &store.ErrorResourceDoesNotExist{}) {
				code = http.StatusNotFound
			}
			err = fmt.Errorf("find plant by id: %w", err)
			logger.Error(err.Error())
			_ = encode(w, r, code, newHttpError(err))
			return
		}

		_ = encode(w, r, http.StatusOK, plant)
	})
}

func handleCreatePlant(plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.LoggerFromCtx(ctx)
		newPlant, problems, err := decodeValid[plants.Plant](r)
		if err != nil {
			err = fmt.Errorf("validation error: %w", err)
			logger.Error(err.Error())
			_ = encode(w, r, http.StatusUnprocessableEntity, newValidationError(err.Error(), problems))
			return
		}

		plant, err := plantStore.Create(ctx, newPlant)
		if err != nil {
			err = fmt.Errorf("create plant: %w", err)
			logger.Error(err.Error())
			_ = encode(w, r, http.StatusInternalServerError, newHttpError(err))
			return
		}

		_ = encode(w, r, http.StatusOK, plant)
	})
}
