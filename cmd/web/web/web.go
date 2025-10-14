package web

import (
	"advent2024/pkg/solver"
	"advent2024/web/middleware"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	logger := middleware.GetLogger(r.Context())

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
		Endpoint: fmt.Sprintf("/api/solve/%s/%s", r.PathValue("day"), r.PathValue("part")),
	}

	if err := template.Execute(w, data); err != nil {
		logger.Printf("Unable to render Upload template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK CI"))
}
