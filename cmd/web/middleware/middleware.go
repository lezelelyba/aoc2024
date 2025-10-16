package middleware

import (
	"advent2024/web/config"
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

const (
	ContextConfig              contextKey = "config"
	ContextKeyLogger           contextKey = "logger"
	ContextKeyRequestID        contextKey = "requestID"
	ContextKeyUploadTemplate   contextKey = "uploadTemplate"
	ContextKeyCallbackTemplate contextKey = "callbackTemplate"
)

type contextKey string

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

func NewLogger(c *config.Config) *log.Logger {
	return log.New(os.Stderr, "advent2024.web ", log.Ldate|log.Ltime|log.LUTC|log.Lmsgprefix)
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			ip := clientIP(r)

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

func GetLogger(r *http.Request) *log.Logger {
	loggerVal := r.Context().Value(ContextKeyLogger)
	logger, ok := loggerVal.(*log.Logger)
	if !ok || logger == nil {
		logger = log.Default()
	}

	return logger
}

func WithTemplate(template *template.Template, key contextKey, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), key, template)
		next(w, r.WithContext(ctx))
	}
}

func RateLimitMiddleware(tokenRate, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(tokenRate), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func WithConfig(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextConfig, cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetConfig(r *http.Request) (*config.Config, bool) {
	cfg, ok := r.Context().Value(ContextConfig).(*config.Config)
	return cfg, ok
}

func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")

	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}
