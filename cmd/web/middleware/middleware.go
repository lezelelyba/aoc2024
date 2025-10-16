package middleware

import (
	"advent2024/web/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func WithTemplate(template *template.Template, key contextKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, template)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
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

func AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logger := GetLogger(r)
			cfg, ok := GetConfig(r)

			if !ok {
				logger.Printf("solve with upload: unable to get config")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Server configuration error"))
				return
			}

			auth := r.Header.Get("Authorization")

			if !strings.HasPrefix(auth, "Bearer ") {
				logger.Printf("unauthorized %s", auth)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil

			})

			if err != nil || !token.Valid {
				logger.Printf("unauthorized: invalid token %s", auth)
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateToken(token string, cfg *config.Config) (bool, error) {
	return true, nil
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

func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
