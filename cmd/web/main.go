package main

import (
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
//	@contact.url	http://localhost

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	// create new mux
	mux := http.NewServeMux()

	// "parse" config

	config := config.Config{Empty: ""}

	// parse templates

	uploadTemplate := template.Must(template.ParseFiles("./templates/upload.tmpl"))

	// create logging middleware

	logger := middleware.NewLogger(&config)

	// web pages
	mux.HandleFunc("GET /", web.ServerStatus)
	mux.HandleFunc("GET /list", web.SolverListing)
	mux.HandleFunc("GET /solve/{day}/{part}", middleware.WithTemplate(uploadTemplate, middleware.ContextKeyUploadTemplate, web.SolveWithUpload))
	mux.HandleFunc("GET /healthcheck", web.HealthCheck)

	// swagger docs
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// api
	mux.HandleFunc("POST /api/solve/{day}/{part}", api.Solve)

	// static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// add logging middleware
	loggedMux := middleware.LoggingMiddleware(logger)(mux)

	// start server
	log.Println("Starting Server on : 8080")
	http.ListenAndServe(":8080", loggedMux)
}
