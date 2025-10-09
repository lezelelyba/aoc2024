package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "advent2024/pkg/d1"
	_ "advent2024/pkg/d2"
	_ "advent2024/pkg/d3"
	_ "advent2024/pkg/d4"
	_ "advent2024/pkg/d5"
	_ "advent2024/pkg/d6"
	_ "advent2024/pkg/d7"

	"advent2024/pkg/solver"

	"github.com/google/uuid"
)

type contextKey string
type Config struct {
	empty string
}

const (
	ContextKeyLogger         contextKey = "logger"
	ContextKeyRequestID      contextKey = "requestID"
	ContextKeyUploadTemplate contextKey = "uploadTemplate"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n

	return n, err
}

func NewLogger(c *Config) *log.Logger {
	return log.New(os.Stderr, "advent2024.web ", log.Ldate|log.Ltime|log.LUTC|log.Lmsgprefix)
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			ip := ClientIP(r)

			requestID := uuid.New().String()

			oldPrefix := logger.Prefix()
			newPrefix := oldPrefix + ": " + requestID + " : "

			prefixedLogger := log.New(logger.Writer(), newPrefix, logger.Flags())

			ctx := context.WithValue(r.Context(), ContextKeyLogger, prefixedLogger)
			ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)

			captureWriter := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(captureWriter, r.WithContext(ctx))

			duration := time.Since(start)

			prefixedLogger.Printf("%s - - \"%s %s\" %d %d duration %s",
				ip,
				r.Method,
				r.URL,
				captureWriter.statusCode,
				captureWriter.bytesWritten,
				duration)
		})
	}
}

func ClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")

	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

func withTemplate(template *template.Template, key contextKey, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), key, template)
		next(w, r.WithContext(ctx))
	}
}

func main() {
	mux := http.NewServeMux()

	// parse config

	config := Config{empty: ""}

	// parse templates

	uploadTemplate := template.Must(template.ParseFiles("./templates/upload.tmpl"))

	// create logging middleware

	logger := NewLogger(&config)

	mux.HandleFunc("GET /", serverStatus)
	mux.HandleFunc("GET /list", solverListing)
	mux.HandleFunc("POST /solve/{day}/{part}", solve)
	mux.HandleFunc("GET /solve/{day}/{part}/upload", withTemplate(uploadTemplate, ContextKeyUploadTemplate, solveWithUpload))
	mux.HandleFunc("GET /healthcheck", healthCheck)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	loggedMux := LoggingMiddleware(logger)(mux)

	log.Println("Starting Server on : 8080")
	http.ListenAndServe(":8080", loggedMux)
}

func serverStatus(w http.ResponseWriter, r *http.Request) {

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

func solverListing(w http.ResponseWriter, r *http.Request) {
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

func getLogger(ctx context.Context) *log.Logger {
	loggerVal := ctx.Value(ContextKeyLogger)
	logger, ok := loggerVal.(*log.Logger)
	if !ok || logger == nil {
		logger = log.Default()
	}

	return logger
}

type SolvePayload struct {
	Input string `json:input`
}

type SolveResult struct {
	Output string `json:output`
}

func solve(w http.ResponseWriter, r *http.Request) {

	logger := getLogger(r.Context())

	day := r.PathValue("day")
	part := r.PathValue("part")

	part_converted, err := strconv.Atoi(part)

	if err != nil {
		logger.Printf("Part is not numerical: %d", part_converted)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Solver for day %s part %s not implemented\n", day, part)))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Printf("Unable to read body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to read body")))
		return
	}

	var p SolvePayload
	err = json.Unmarshal(body, &p)

	if err != nil {
		logger.Printf("Unable to unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to unmarshal request")))
		return
	}

	decoded_body, err := base64.StdEncoding.DecodeString(string(p.Input))

	if err != nil {
		logger.Printf("Unable to decode Base64")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to read body: Invalid Base64 encoding")))
		return
	}

	solver, ok := solver.New(day)

	if !ok {
		logger.Printf("Unable to find solver for day %s", day)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Solver for day %s not implemented\n", day)))
		return
	}

	err = solver.Init(strings.NewReader(string(decoded_body)))

	if err != nil {
		logger.Printf("Unable to intialize Solver for day %s", day)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Unable to initialize Solver for day %s\n", day)))
		return
	}

	result, err := solver.Solve(part_converted)

	if err != nil {
		logger.Printf("Unable to solve problem for day %s part %s", day, part)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to solve for day %s\n", day)))
		return
	}

	b, err := json.Marshal(SolveResult{Output: result})
	if err != nil {
		logger.Printf("Unable to marshal result")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func solveWithUpload(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(r.Context())

	templateVal := r.Context().Value(ContextKeyUploadTemplate)
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
		Endpoint: fmt.Sprintf("/solve/%s/%s", r.PathValue("day"), r.PathValue("part")),
	}

	if err := template.Execute(w, data); err != nil {
		logger.Printf("Unable to render Upload template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK CI Test 2nd"))
}
