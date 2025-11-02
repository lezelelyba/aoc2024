// Solver for Advent of Code 2024
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "advent2024/pkg/d1"
	_ "advent2024/pkg/d10"
	_ "advent2024/pkg/d11"
	_ "advent2024/pkg/d2"
	_ "advent2024/pkg/d3"
	_ "advent2024/pkg/d4"
	_ "advent2024/pkg/d5"
	_ "advent2024/pkg/d6"
	_ "advent2024/pkg/d7"
	_ "advent2024/pkg/d8"
	_ "advent2024/pkg/d9"

	_ "advent2024/web/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"advent2024/web/api"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/webhandlers"
)

var Version string = "dev"

//	@title			Advent of Code 2024 Solver API
//	@version		2.0
//	@description	Solver for AoC 2024 written in Go

//	@contact.name	None

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

// @securitydefinitions.oauth2.accessCode	OAuth2AccessCode
// @authorizationURL						https://github.com/login/oauth/authorize
// @tokenURL								http://localhost:8080/oauth/github/token
// @scope.read								Grants read access
// @description							GitHub OAuth
func main() {

	// parse and validate config.

	cfg, errors := config.LoadConfig()
	valid, validationErrors := cfg.ValidateConfig()

	if len(errors) != 0 {
		for _, e := range errors {
			log.Println(e)
		}
	}

	if len(validationErrors) != 0 {
		for _, e := range validationErrors {
			log.Println(e)
		}
	}

	if !valid || len(errors) > 0 || len(validationErrors) > 0 {
		os.Exit(1)
	}

	cfg.Version = Version

	// parse templates

	funcMap := template.FuncMap{
		"fieldNames": webhandlers.FieldNames,
		"getField":   webhandlers.GetField,
	}

	// base layout
	layoutTemplate := template.Must(template.New("").Funcs(funcMap).ParseGlob("./templates/layouts/*.tmpl"))

	// index page
	indexTemplate := template.Must(layoutTemplate.Clone())
	template.Must(indexTemplate.ParseFiles("./templates/pages/index.tmpl"))

	// doc page
	docsTemplate := template.Must(layoutTemplate.Clone())
	template.Must(docsTemplate.ParseFiles("./templates/pages/docs.tmpl"))

	// callback page
	callbackTemplate := template.Must(template.ParseFiles("./templates/pages/callback.tmpl"))

	// create logging middleware

	logger := middleware.NewLogger(&cfg)

	// create http muxes
	webMux := http.NewServeMux()
	apiMux := http.NewServeMux()
	apiUnsecuredMux := http.NewServeMux()
	globalMux := http.NewServeMux()

	// web pages
	webMux.Handle("GET /", middleware.WithTemplate(indexTemplate)(http.HandlerFunc(webhandlers.ServerIndex)))
	webMux.Handle("GET /docs", middleware.WithTemplate(docsTemplate)(http.HandlerFunc(webhandlers.ServerDocs)))
	webMux.HandleFunc("GET /list", webhandlers.SolverListing)
	webMux.HandleFunc("GET /healthcheck", webhandlers.HealthCheck)

	// oauth
	if cfg.OAuth {
		webMux.Handle("GET /callback/{provider}", middleware.WithTemplate(callbackTemplate)(http.HandlerFunc(webhandlers.OAuthCallback)))
	}

	// swagger docs
	webMux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// static files
	fs := http.FileServer(http.Dir("static"))
	webMux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// godoc
	godocs := http.FileServer(http.Dir("godocs"))
	webMux.Handle("GET /godocs/", http.StripPrefix("/godocs/", godocs))

	// api
	apiMux.HandleFunc("GET /solvers", api.SolverListing)
	apiMux.HandleFunc("POST /solvers/{day}/{part}", api.Solve)

	// add api rate limiter
	apiHandler := middleware.RateLimitMiddleware(cfg.APIRate, cfg.APIBurst)(apiMux)

	// add authentication if enabled
	if cfg.OAuth {
		apiHandler = middleware.AuthenticationMiddleware()(apiHandler)

		// token exchange token
		apiUnsecuredMux.HandleFunc("POST /access_token", api.OAuthCodeExchange)
		apiUnsecuredHandler := middleware.RateLimitMiddleware(cfg.APIRate, cfg.APIBurst)(apiUnsecuredMux)
		globalMux.Handle("/api/public/", http.StripPrefix("/api/public", apiUnsecuredHandler))
	}

	// combine muxes
	globalMux.Handle("/api/", http.StripPrefix("/api", apiHandler))
	globalMux.Handle("/", webMux)

	// add middlewares
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
