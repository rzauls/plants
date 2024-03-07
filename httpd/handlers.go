package httpd

import (
	"fmt"
	"net/http"
)

// ListPlants godoc
// @Summary      List all plants
// @Description  Get all plants
// @Tags         plants
// @Accept       json
// @Produce      json
// @Success      200  {array}   plants.Plant
// @Failure      400  {object}  httpError
// @Failure      404  {object}  httpError
// @Failure      500  {object}  httpError
// @Router       /plants [get]
func (s *httpService) handleListPlants(w http.ResponseWriter, r *http.Request) {
	plants, err := s.store.List()
	if err != nil {
		errorMsg := fmt.Errorf("retrieve all plants: %w", err)
		s.logger.Error(errorMsg.Error())
		s.response(w, http.StatusInternalServerError, newHttpError(errorMsg))
		return
	}

	s.response(w, http.StatusOK, plants)
}

// GetPlant godoc
// @Summary      Get plant
// @Description  Get plant by id
// @Tags         plants
// @Accept       json
// @Produce      json
// @Success      200  {object}  plants.Plant
// @Failure      400  {object}  httpError
// @Failure      404  {object}  httpError
// @Failure      500  {object}  httpError
// @Router       /plants/{id} [get]
func (s *httpService) handleGetPlant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	plant, err := s.store.Find(id)
	if err != nil {
		errorMsg := fmt.Errorf("find plant by id: %w", err)
		s.logger.Error(errorMsg.Error())
		// NOTE: if store has typed errors, we can change http response codes here
		s.response(w, http.StatusInternalServerError, newHttpError(errorMsg))
		return
	}

	s.response(w, http.StatusOK, plant)
}
