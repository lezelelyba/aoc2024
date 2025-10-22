package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "advent2024/pkg/d1"
	_ "advent2024/pkg/d10"
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

var Version string = "dev"

//	@title			Advent of Code 2024 Solver API
//	@version		1.0
//	@description	Solver for AoC 2024 written in Go

//	@contact.name	None

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

//	@securitydefinitions.oauth2.accessCode	OAuth2AccessCode
//	@authorizationURL						https://github.com/login/oauth/authorize
//	@tokenURL								http://localhost:8080/oauth/github/token
//	@scope.read								Grants read access
//	@description							GitHub OAuth

func main() {
	// create new mux
	webMux := http.NewServeMux()
	apiMux := http.NewServeMux()
	globalMux := http.NewServeMux()

	// parse config

	cfg, errors := config.LoadConfig()

	if len(errors) != 0 {
		for _, e := range errors {
			log.Println(e)
		}
	}

	cfg.Version = Version

	// parse templates

	funcMap := template.FuncMap{
		"fieldNames": web.FieldNames,
		"getField":   web.GetField,
	}

	indexTemplate := template.Must(template.New("").Funcs(funcMap).ParseFiles("./templates/index.tmpl"))
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

	logger := middleware.NewLogger(&cfg)

	// web pages
	webMux.Handle("GET /", middleware.WithTemplate(indexTemplate)(http.HandlerFunc(web.ServerIndex)))
	webMux.HandleFunc("GET /list", web.SolverListing)
	webMux.HandleFunc("GET /healthcheck", web.HealthCheck)

	// oauth
	if cfg.OAuth {
		webMux.Handle("GET /callback/{provider}", middleware.WithTemplate(callbackTemplate)(http.HandlerFunc(web.OAuthCallback)))
		webMux.HandleFunc("POST /oauth/{provider}/token", web.OAuthHandler)
	}

	// swagger docs
	webMux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// static files
	fs := http.FileServer(http.Dir("static"))
	webMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// api
	apiMux.HandleFunc("GET /solvers", api.SolverListing)
	apiMux.HandleFunc("POST /solvers/{day}/{part}", api.Solve)

	// add api rate limiter
	apiHandler := middleware.RateLimitMiddleware(cfg.APIRate, cfg.APIBurst)(apiMux)

	// add authentication if enabled
	if cfg.OAuth {
		apiHandler = middleware.AuthenticationMiddleware()(apiHandler)
	}

	// combine muxes
	globalMux.Handle("/api/", http.StripPrefix("/api", apiHandler))
	globalMux.Handle("/", webMux)

	// add logging middleware
	var finalMux http.Handler = globalMux
	finalMux = middleware.RecoveryMiddleware()(finalMux)
	finalMux = middleware.LoggingMiddleware(logger)(finalMux)
	finalMux = middleware.WithConfig(&cfg)(finalMux)

	// start server
	addr := fmt.Sprintf(":%d", cfg.Port)

	if !cfg.EnableTLS {
		log.Printf("Starting Server on : %d\n", cfg.Port)
		log.Fatal(http.ListenAndServe(addr, finalMux))
	} else {
		log.Printf("Starting TLS Server on : %d\n", cfg.Port)
		log.Fatal(http.ListenAndServeTLS(addr, cfg.CertFile, cfg.KeyFile, finalMux))
	}
}
