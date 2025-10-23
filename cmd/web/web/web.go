// Web Page Handlers
package web

import (
	"advent2024/pkg/solver"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// Server Status
// including hostname and list of registered solvers
func ServerStatus(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	registeredKeys := solver.ListRegistryItemsWithCtx()
	var keyNames []string

	for _, i := range registeredKeys {
		keyNames = append(keyNames, i.Name)
	}
	registeredKeysStr := strings.Join(keyNames, " ")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownHostname"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server " + hostname + " is up and running\n" + "Registered days: " + registeredKeysStr + "\n"))
}

// Main Page
// Renders Main page from template
// TODO: refactor to be usable for all simple pages
func ServerIndex(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	var page = "index"

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("configuration error: %s: unable to get config", page)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyTemplates)
	tmpl, ok := templateVal.(*template.Template)

	ok = ok && tmpl != nil

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to find %s template", page)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownHostname"
	}

	registryItems := solver.ListRegistryItemsWithCtx()

	data := struct {
		Hostname      string
		Config        *config.Config
		RegistryItems []solver.RegistryItemPublic
	}{
		Hostname:      hostname,
		RegistryItems: registryItems,
		Config:        cfg,
	}

	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render %s template %v", page, err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

// Docs page Handler
// Renders page from template
func ServerDocs(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	page := "docs"

	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("configuration error: %s: unable to get config", page)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyTemplates)
	tmpl, ok := templateVal.(*template.Template)

	ok = ok && tmpl != nil

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to find %s template", page)
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "UnknownHostname"
	}

	registryItems := solver.ListRegistryItemsWithCtx()

	data := struct {
		Hostname      string
		Config        *config.Config
		RegistryItems []solver.RegistryItemPublic
	}{
		Hostname:      hostname,
		RegistryItems: registryItems,
		Config:        cfg,
	}

	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render %s template %v", page, err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

// Shows registered Solvers
func SolverListing(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	registryItems := solver.ListRegistryItemsWithCtx()

	b, err := json.Marshal(registryItems)

	rc := http.StatusInternalServerError
	errMsg := "unable to marchal registered items solvers"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// Handles callback redirects from OAuth providers
// Pulls data about OAuth provider from config
// Renders page to send user to Handler URL
func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	providerName := r.PathValue("provider")
	_, exists := config.OAuthProviders[providerName]

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)
	if weberrors.HandleError(w, logger, weberrors.OkToError(exists), rc, errMsg) != nil {
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyTemplates)
	template, ok := templateVal.(*template.Template)

	ok = ok && template != nil

	rc = http.StatusInternalServerError
	errMsg = "unable to find callback template"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	provider := config.OAuthProviders[providerName]

	data := struct {
		Endpoint string
	}{
		Endpoint: provider.TokenEndpoint(),
	}

	err := template.Execute(w, data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render callback template %v", err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

// Handles code<>token exchange with the OAuth Providers
// TODO: implement more providers, curently only github
// gets OAuth code from client, exchanges code for token with provider and then generates JWT token for client
func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	// get provider name
	providerName := r.PathValue("provider")
	_, exists := config.OAuthProviders[providerName]

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)
	if weberrors.HandleError(w, logger, weberrors.OkToError(exists), rc, errMsg) != nil {
		return
	}

	// get provider configuration
	query := r.URL.Query()
	provider := config.OAuthProviders[providerName]

	switch provider.Name {
	case "github":

		// code is missing
		rc = http.StatusBadRequest
		errMsg = fmt.Sprintf("unable to exchange code for token with %s: code is missing", provider.Name)
		ok := query.Get("code") != ""
		if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
			return
		}

		// exchange code
		token, err := exchangeCodeForToken(&provider, query.Get("code"))

		// error => local error
		rc = http.StatusInternalServerError
		errMsg = fmt.Sprintf("unable to exchange code for token with %s: %v", provider.Name, err)
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		// error nil, but error response from provider
		if token.IsError() {
			rc = http.StatusBadRequest
			errMsg = fmt.Sprintf("unable to exchange code for token with %s: %v", provider.Name, err)
			if weberrors.HandleError(w, logger, token, rc, errMsg) != nil {
				return
			}
		}

		// prepare response for client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// get token and generate JWT token
		tokenStr, _ := token.Token()
		jwtToken, err := middleware.GenerateJWT(provider.Name, tokenStr, []byte(config.JWTSecret), config.JWTTokenValidity)

		// unable to generate token
		rc = http.StatusInternalServerError
		errMsg = "unable to create jwt token"
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		// write response to client
		json.NewEncoder(w).Encode(map[string]string{"access_token": jwtToken})

		return
	}
}

// Helper function for the actual code<>token exchange
func exchangeCodeForToken(provider *config.OAuthProvider, code string) (middleware.OAuthResponse, error) {
	// no provider
	if provider == nil {
		return nil, fmt.Errorf("unable to find empty provider")
	}

	// know providres
	switch provider.Name {
	case "github":
		// extract required information from client request
		data := url.Values{}
		data.Set("client_id", provider.ClientId)
		data.Set("client_secret", provider.ClientSecret)
		data.Set("code", code)
		// TODO: send redirect_uri to github

		// create request to provider
		req, err := http.NewRequest(
			"POST",
			provider.TokenURL,
			strings.NewReader(data.Encode()))

		if err != nil {
			return nil, fmt.Errorf("unable to create OAuth request")
		}

		// set headers and make the request
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)

		// err from client
		if err != nil {
			return nil, fmt.Errorf("unable to exchange code for token with %s: %v", provider.Name, err)
		}

		// work only with non-nil response
		defer resp.Body.Close()

		// decode token
		var token middleware.OAuthGithubReply

		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal %s response", provider.Name)
		}

		return token, nil
	// unknown provider
	default:
		return nil, fmt.Errorf("unable to find provider %s", provider.Name)
	}
}

// Healthcheck page Handler
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// helper functions for templates, extract field names from struct
func FieldNames(v interface{}) []string {
	r := reflect.ValueOf(v)
	t := r.Type()
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath == "" { // exported
			fields = append(fields, f.Name)
		}
	}
	return fields
}

// helper functions for templates, extracts value of field from struct
func GetField(v interface{}, name string) interface{} {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Ptr {
		r = r.Elem() // dereference pointer
	}
	if r.Kind() == reflect.Struct {
		f := r.FieldByName(name)
		if f.IsValid() {
			return f.Interface()
		}
	}
	return nil
}
