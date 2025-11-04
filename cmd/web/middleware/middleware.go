// Collection of Middlewares
package middleware

import (
	"advent2024/web/config"
	"advent2024/web/weberrors"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// Context Keys
const (
	ContextConfig       contextKey = "config"
	ContextKeyLogger    contextKey = "logger"
	ContextKeyRequestID contextKey = "requestID"
	ContextKeyTemplates contextKey = "uploadTemplate"
)

// Context key type
type contextKey string

// Logger to capture status code and written bytes
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// Writes header and saves the code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Writes the response and saves the amount of bytes
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n

	return n, err
}

// Creates new longer based on the configuration
// TODO: add more options
func NewLogger(c *config.Config) *log.Logger {
	return log.New(os.Stderr, "advent2024.web ", log.Ldate|log.Ltime|log.LUTC|log.Lmsgprefix)
}

// Logs return code, written bytes and duration of the request
// Generates unique request ID
// TODO: separate middleware with stat tracking?
// Logs are prefixed with timestamp and project name
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// mark start time
			start := time.Now()

			// get user IP
			ip := clientIP(r)

			// generate new request id
			requestID := uuid.New().String()

			// add request ID to the logger prefix
			oldPrefix := logger.Prefix()
			newPrefix := oldPrefix + ": " + requestID + " : "
			prefixedLogger := log.New(logger.Writer(), newPrefix, logger.Flags())

			// insert logger into context
			ctx := context.WithValue(r.Context(), ContextKeyLogger, prefixedLogger)
			ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)

			// insert capture writter into the handler chain
			captureWriter := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(captureWriter, r.WithContext(ctx))

			// capture request duration
			duration := time.Since(start)

			// write response
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

// Gets logger from context or provides a default one
func GetLogger(r *http.Request) *log.Logger {
	loggerVal := r.Context().Value(ContextKeyLogger)
	logger, ok := loggerVal.(*log.Logger)
	if !ok || logger == nil {
		logger = log.Default()
	}

	return logger
}

// Specifies templates to be used by the Handler
func WithTemplate(template *template.Template) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextKeyTemplates, template)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Rate Limits the request to Handler
func RateLimitMiddleware(tokenRate, burst int) func(http.Handler) http.Handler {
	// create new limiter based on params
	limiter := rate.NewLimiter(rate.Limit(tokenRate), burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := GetLogger(r)
			// if allowed serve, otherwise reply with TooManyRequests
			if limiter.Allow() {
				next.ServeHTTP(w, r)
			} else {
				rc := http.StatusTooManyRequests
				errMsg := "too many requests"
				_ = weberrors.HandleError(w, logger, errors.New(errMsg), rc, errMsg)

				return
			}
		})
	}
}

// Checks user authentication
func AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var rc int
			var errMsg string

			// get config and logger
			logger := GetLogger(r)
			cfg, ok := GetConfig(r)

			rc = http.StatusInternalServerError
			errMsg = "configuration error: unable to get config"
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			// TODO: add other forms of authentication

			// Check if bearer token is present
			auth := r.Header.Get("Authorization")
			ok = strings.HasPrefix(auth, "Bearer ")

			// if missing or wrong format => unauthorized
			rc = http.StatusUnauthorized
			errMsg = fmt.Sprintf("unauthorized")
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			// parse and validate token
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			token, _ := ParseToken(tokenStr, []byte(cfg.JWTSecret))

			ok = TokenValid(token)

			// if token is invalid or expired => unauthorized
			rc = http.StatusUnauthorized
			errMsg = fmt.Sprintf("unauthorized: invalid token")
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORS
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Injects configuration to Handler chain
func WithConfig(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextConfig, cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Gets Configuration from context
// TODO: provide default configuration if missing? Similar to logger?
func GetConfig(r *http.Request) (*config.Config, bool) {
	cfg, ok := r.Context().Value(ContextConfig).(*config.Config)
	if cfg == nil {
		return cfg, false
	}
	return cfg, ok
}

// Recovers from panics within the Handler chain
func RecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					var rc int
					var errMsg string

					logger := GetLogger(r)

					rc = http.StatusInternalServerError
					errMsg = "internal server error: recovered"
					if weberrors.HandleError(w, logger, weberrors.OkToError(false), rc, errMsg) != nil {
						return
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Extracts client IP from the request
// X-Real-Ip header > X-Forwarded-For header > src IP of request
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

// Chains Middlewares
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
