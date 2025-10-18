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

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

// @securitydefinitions.oauth2.accessCode	OAuth2AccessCode
// @authorizationURL						https://github.com/login/oauth/authorize
// @tokenURL								https://github.com/login/oauth/access_token
// @scope.read								Grants read access
// @description							GitHub OAuth
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

	funcMap := template.FuncMap{
		"fieldNames": web.FieldNames,
		"getField":   web.GetField,
	}

	indexTemplate := template.Must(template.New("").Funcs(funcMap).ParseFiles("./templates/index.tmpl"))
	// uploadTemplate := template.Must(template.ParseFiles("./templates/upload.tmpl"))
	callbackTemplate := template.Must(template.ParseFiles("./templates/callback.tmpl"))

	// TODO: load all common templates
	// TODO: create template for each page, including the common templates
	// use those templates to render each page
	//
	// cannot load all templates at once and then reference just one of them
	// as for the common blocks, the last template loaded which defines that block
	// will be use
	//
	// templates := template.Must(template.ParseGlob("./templates/*.tmpl"))

	// create logging middleware

	logger := middleware.NewLogger(&config)

	// web pages
	webMux.Handle("GET /",
		middleware.Chain(
			http.HandlerFunc(web.ServerIndex),
			middleware.WithConfig(&config),
			middleware.WithTemplate(indexTemplate, middleware.ContextKeyIndexTemplate)))
	webMux.HandleFunc("GET /list", web.SolverListing)
	webMux.HandleFunc("GET /healthcheck", web.HealthCheck)
	// webMux.Handle("GET /solve/{day}/{part}",
	// 	middleware.Chain(
	// 		http.HandlerFunc(web.SolveWithUpload),
	// 		middleware.WithConfig(&config),
	// 		middleware.WithTemplate(uploadTemplate, middleware.ContextKeyUploadTemplate)))

	// oauth
	if config.OAuth {
		webMux.Handle("GET /callback/{provider}",
			middleware.Chain(
				http.HandlerFunc(web.OAuthCallback),
				middleware.WithConfig(&config),
				middleware.WithTemplate(callbackTemplate, middleware.ContextKeyCallbackTemplate)))
		webMux.Handle("POST /oauth/{provider}/token",
			middleware.Chain(
				http.HandlerFunc(web.OAuthHandler),
				middleware.WithConfig(&config)))
	}

	// swagger docs
	webMux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// static files
	fs := http.FileServer(http.Dir("static"))
	webMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// api
	apiMux.HandleFunc("GET /solvers", api.SolverListing)
	apiMux.HandleFunc("POST /solvers/{day}/{part}", api.Solve)

	// combine muxes
	apiHandler := middleware.RateLimitMiddleware(config.APIRate, config.APIBurst)(apiMux)
	if config.OAuth {
		apiHandler = middleware.Chain(
			apiHandler,
			middleware.WithConfig(&config),
			middleware.AuthenticationMiddleware())
	}

	globalMux.Handle("/api/", http.StripPrefix("/api", apiHandler))
	globalMux.Handle("/", webMux)

	// add logging middleware
	finalMux := middleware.LoggingMiddleware(logger)(globalMux)

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
