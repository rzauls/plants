package httpd

import (
	"net/http"

	"github.com/swaggo/http-swagger/v2"

	_ "plants/docs"
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
	host := "localhost:8080"
	mux := http.NewServeMux()
    mux.handleFunc("/docs/swagger.json", 
    http.Handle("/",
    http.FileServer(
        http.File("../docs/swagger.json"),
    )
)
	mux.HandleFunc("/docs/*", httpSwagger.Handler(
		httpSwagger.URL(host+"/docs/swagger.json"),
	),
	)
	mux.HandleFunc(http.MethodGet+" /plants/", handleListPlants)
	mux.HandleFunc(http.MethodGet+" /plants/{id}/", handleGetPlant)

	http.ListenAndServe(host, mux)
}
