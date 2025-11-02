// Web Page Handlers
package webhandlers

import (
	"advent2024/pkg/solver"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// Handles server status page.
// Output including hostname and list of registered solvers
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

// Handles main page
// Renders page from template, with information pulled from configuration
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

// Handles documentation page
// Renders page from template, with information pulled from configuration
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

// Handles display of registered solvers
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
// Renders page from template, with information pulled from configuration
func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	providerName := r.PathValue("provider")
	_, exists := cfg.OAuthProviders[providerName]

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

	provider := cfg.OAuthProviders[providerName]

	data := struct {
		Provider config.OAuthProvider
	}{
		Provider: provider,
	}

	err := template.Execute(w, data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render callback template %v", err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

// Handles healthcheck page
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// extract field names from struct
// helper functions for templates
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

// extracts value of field from struct
// helper functions for templates
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
