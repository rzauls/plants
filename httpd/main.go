package httpd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "plants/docs"
	"plants/plants"
	"plants/store"
)

// @title           Swagger Plant API
// @version         1.0
// @description     This is a sample server plant server with semi-auto generated swagger docs

// @contact.name   Rihards Zauls
// @contact.email  rihards.zauls@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func Run() {
	// NOTE: realistically this wouldnt be an in-memory array,
	// but a DB implementation of store.Store interface

	// TODO: add context for stopping everything
	s := store.NewMemoryStore([]plants.Plant{})
	logger := slog.Default()

	service := httpService{
		store:  s,
		logger: logger,
	}

	host := "localhost:8080"
	mux := http.NewServeMux()
	// mux.handleFunc("/docs/swagger.json",
	// http.Handle("/",
	// http.FileServer(
	//     http.File("../docs/swagger.json"),
	// )
	// )
	mux.HandleFunc("/docs/*", httpSwagger.Handler(
		httpSwagger.URL(host+"/docs/swagger.json"),
	),
	)
	mux.HandleFunc(http.MethodGet+" /plants/", service.handleListPlants)
	mux.HandleFunc(http.MethodGet+" /plants/{id}/", service.handleGetPlant)

	logger.Info(fmt.Sprintf("Listening to requests on %s\n", host))

	//set up middleware chain
	chain := logRequestsMiddleware(mux)
	http.ListenAndServe(host, chain)
}

type httpService struct {
	store  store.Store
	logger *slog.Logger
}

func (s *httpService) response(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.logger.Error(fmt.Errorf("marshal response payload: %w", err).Error())
	}
}

func logRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Default().Info(fmt.Sprintf(
			"%s %s",
			r.Method,
			r.URL.String(),
		))
		// NOTE: could wrap this with another writer so we can log the response code aswell, but since we log errors this isnt strictly useful
		next.ServeHTTP(w, r)
	})
}
