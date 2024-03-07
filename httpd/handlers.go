package httpd

import (
	"errors"
	"fmt"
	"net/http"
	"plants/plants"
)

// ListPlants godoc
//
//	@Summary		List all plants
//	@Description	Get all plants
//	@Tags			plants
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		plants.Plant
//	@Failure		400	{object}	httpError
//	@Failure		404	{object}	httpError
//	@Failure		500	{object}	httpError
//	@Router			/plants [get]
func (s *httpService) handleListPlants(w http.ResponseWriter, r *http.Request) {
	plts, err := s.store.List()
	if err != nil {
		s.responseError(w, http.StatusInternalServerError, fmt.Errorf("retrieve all plants: %w", err))
		return
	}

	if len(plts) == 0 {
		plts = make([]plants.Plant, 0)
	}

	s.response(w, http.StatusOK, plts)
}

// GetPlant godoc
//
//	@Summary		Get plant
//	@Description	Get plant by id
//	@Tags			plants
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	plants.Plant
//	@Failure		400	{object}	httpError
//	@Failure		404	{object}	httpError
//	@Failure		500	{object}	httpError
//	@Router			/plants/{id} [get]
func (s *httpService) handleGetPlant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.responseError(w, http.StatusUnprocessableEntity, errors.New("id is required in path parameters"))
		return
	}

	plant, err := s.store.Find(id)
	if err != nil {
		s.responseError(w, http.StatusInternalServerError, fmt.Errorf("find plant by id: %w", err))
		return
	}

	if plant == nil {
		s.responseError(w, http.StatusNotFound, fmt.Errorf("plant with ID '%s' does not exist", id))
		return
	}

	s.response(w, http.StatusOK, plant)
}
