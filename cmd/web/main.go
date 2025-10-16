package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "advent2024/pkg/d1"
	_ "advent2024/pkg/d2"
	_ "advent2024/pkg/d3"
	_ "advent2024/pkg/d4"
	_ "advent2024/pkg/d5"
	_ "advent2024/pkg/d6"
	_ "advent2024/pkg/d7"
	_ "advent2024/pkg/d8"
	_ "advent2024/pkg/d9"

	httpSwagger "github.com/swaggo/http-swagger"

	"advent2024/web/api"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/web"

	_ "advent2024/web/docs"
)

//	@title			Advent of Code 2024 Solver API
//	@version		1.0
//	@description	Solver for AoC 2024 written in Go

//	@contact.name	None

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// create new mux
	webMux := http.NewServeMux()
	apiMux := http.NewServeMux()
	globalMux := http.NewServeMux()

	// parse config

	config, errors := config.LoadConfig()

	if len(errors) != 0 {
		for _, e := range errors {
			log.Println(e)
		}
	}

	// parse templates

	uploadTemplate := template.Must(template.ParseFiles("./templates/upload.tmpl"))
	callbackTemplate := template.Must(template.ParseFiles("./templates/callback.tmpl"))

	// create logging middleware

	logger := middleware.NewLogger(&config)

	// web pages
	webMux.HandleFunc("GET /", web.ServerStatus)
	webMux.HandleFunc("GET /list", web.SolverListing)
	webMux.HandleFunc("GET /solve/{day}/{part}", middleware.WithTemplate(uploadTemplate, middleware.ContextKeyUploadTemplate, web.SolveWithUpload))
	webMux.HandleFunc("GET /healthcheck", web.HealthCheck)

	// oauth
	webMux.HandleFunc("GET /callback", middleware.WithTemplate(callbackTemplate, middleware.ContextKeyCallbackTemplate, web.OAuthCallback))
	webMux.HandleFunc("POST /oauth/{provider}/token", web.OAuthHandler)

	// swagger docs
	webMux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// static files
	fs := http.FileServer(http.Dir("static"))
	webMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// api
	apiMux.HandleFunc("GET /solvers", api.SolverListing)
	apiMux.HandleFunc("POST /solvers/{day}/{part}", api.Solve)

	// combine muxes
	globalMux.Handle("/api/", http.StripPrefix("/api", middleware.RateLimitMiddleware(config.APIRate, config.APIBurst)(apiMux)))
	globalMux.Handle("/", webMux)

	// add logging middleware
	loggedMux := middleware.LoggingMiddleware(logger)(globalMux)

	// add config middleware
	finalMux := middleware.WithConfig(&config)(loggedMux)

	// start server
	addr := fmt.Sprintf(":%d", config.Port)

	if !config.EnableTLS {
		log.Printf("Starting Server on : %d\n", config.Port)
		log.Fatal(http.ListenAndServe(addr, finalMux))
	} else {
		log.Printf("Starting TLS Server on : %d\n", config.Port)
		log.Fatal(http.ListenAndServeTLS(addr, config.CertFile, config.KeyFile, finalMux))
	}
}
