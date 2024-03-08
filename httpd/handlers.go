package httpd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"plants/plants"
	"plants/store"
)

func handleHealth(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if all dependencies are accessible and ready for connections
		encode(w, r, http.StatusOK, "ok")
	})
}

func handleListPlants(logger *slog.Logger, plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plts, err := plantStore.List()
		if err != nil {
			err = fmt.Errorf("retrieve all plants: %w", err)
			logger.Error(err.Error())
			encode(w, r, http.StatusInternalServerError, newHttpError(err))
			return
		}

		if len(plts) == 0 {
			plts = make([]plants.Plant, 0)
		}

		encode(w, r, http.StatusOK, plts)
	})
}

func handleGetPlant(logger *slog.Logger, plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			err := errors.New("id is required in path parameters")
			logger.Error(err.Error())
			encode(w, r, http.StatusUnprocessableEntity, newHttpError(err))
			return
		}

		plant, err := plantStore.Find(id)
		if err != nil {
			err = fmt.Errorf("find plant by id: %w", err)
			logger.Error(err.Error())
			encode(w, r, http.StatusInternalServerError, newHttpError(err))
			return
		}

		if plant == nil {
			err = fmt.Errorf("plant with ID '%s' does not exist", id)
			logger.Error(err.Error())
			encode(w, r, http.StatusNotFound, newHttpError(err))
			return
		}

		encode(w, r, http.StatusOK, plant)
	})
}

func handleCreatePlant(logger *slog.Logger, plantStore store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPlant, problems, err := decodeValid[plants.Plant](r)
		if err != nil {
			err = fmt.Errorf("validation error: %w", err)
			logger.Error(err.Error())
			encode(w, r, http.StatusUnprocessableEntity, newValidationError(err.Error(), problems))
			return
		}

		plant, err := plantStore.Create(newPlant)
		if err != nil {
			err = fmt.Errorf("create plant: %w", err)
			logger.Error(err.Error())
			encode(w, r, http.StatusInternalServerError, newHttpError(err))
			return
		}

		encode(w, r, http.StatusOK, plant)
	})
}
