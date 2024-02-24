package httpd

import (
	"fmt"
	"net/http"
    "plants/plants"
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
func handleListPlants(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "got all plants\n")
    all := []plants.Plant{}
    fmt.Printf("plants: %v", all)

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
func handleGetPlant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "got plant with id=%v\n", id)
}
