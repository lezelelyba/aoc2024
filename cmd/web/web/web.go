package web

import (
	"advent2024/pkg/solver"
	"advent2024/web/middleware"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

func ServerStatus(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	registered_keys := solver.ListRegister()
	registered_keys_string := strings.Join(registered_keys, " ")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownHostname"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server " + hostname + " is up and running\n" + "Registered days: " + registered_keys_string + "\n"))
}

func SolverListing(w http.ResponseWriter, r *http.Request) {
	registered_keys := solver.ListRegister()

	type registeredKeys struct {
		Keys []string
	}

	b, err := json.Marshal(registeredKeys{Keys: registered_keys})
	if err != nil {
		log.Printf("Unable to marshal registered solvers")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func SolveWithUpload(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	templateVal := r.Context().Value(middleware.ContextKeyUploadTemplate)
	template, ok := templateVal.(*template.Template)

	if !ok || template == nil {
		logger.Printf("Unable to find Upload template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	data := struct {
		Title    string
		Endpoint string
	}{
		Title:    fmt.Sprintf("Upload Page for day %s", r.PathValue("day")),
		Endpoint: fmt.Sprintf("/api/solvers/%s/%s", r.PathValue("day"), r.PathValue("part")),
	}

	if err := template.Execute(w, data); err != nil {
		logger.Printf("Unable to render Upload template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	templateVal := r.Context().Value(middleware.ContextKeyCallbackTemplate)
	template, ok := templateVal.(*template.Template)

	if !ok || template == nil {
		logger.Printf("Unable to find Callback template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	if err := template.Execute(w, ""); err != nil {
		logger.Printf("Unable to render Callback template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	if !ok {
		logger.Printf("unable to get config")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server configuration error"))
		return
	}

	providerName := r.PathValue("provider")
	if _, exists := config.OAuthProviders[providerName]; !exists {
		logger.Printf("unable to find oauth provider %s", providerName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown OAuth provider"))
		return
	}

	query := r.URL.Query()
	provider := config.OAuthProviders[providerName]

	switch provider.Name {
	case "github":
		data := url.Values{}
		data.Set("client_id", provider.ClientId)
		data.Set("client_secret", provider.ClientSecret)
		data.Set("code", query.Get("code"))
		// TODO: redirect_uri

		req, err := http.NewRequest(
			"POST",
			provider.URL,
			strings.NewReader(data.Encode()))

		if err != nil {
			logger.Printf("unable to create OAuth Request")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if err != nil {
			logger.Printf("unable to resolve token with %s", provider.Name)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, resp.Body)

		return
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
