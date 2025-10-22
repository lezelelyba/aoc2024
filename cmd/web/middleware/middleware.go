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

const (
	ContextConfig       contextKey = "config"
	ContextKeyLogger    contextKey = "logger"
	ContextKeyRequestID contextKey = "requestID"
	ContextKeyTemplates contextKey = "uploadTemplate"
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

			// print start of the request
			// prefixedLogger.Printf("%s - S \"%s %s\"",
			// 	ip,
			// 	r.Method,
			// 	r.URL)

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

func WithTemplate(template *template.Template) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextKeyTemplates, template)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RateLimitMiddleware(tokenRate, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(tokenRate), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := GetLogger(r)

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

func AuthenticationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var rc int
			var errMsg string

			logger := GetLogger(r)
			cfg, ok := GetConfig(r)

			rc = http.StatusInternalServerError
			errMsg = "configuration error: unable to get config"
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			auth := r.Header.Get("Authorization")
			ok = strings.HasPrefix(auth, "Bearer ")

			rc = http.StatusUnauthorized
			errMsg = fmt.Sprintf("unauthorized")
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			token, _ := ParseToken(tokenStr, []byte(cfg.JWTSecret))

			ok = TokenValid(token)

			rc = http.StatusUnauthorized
			errMsg = fmt.Sprintf("unauthorized: invalid token")
			if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
				return
			}

			next.ServeHTTP(w, r)
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
	if cfg == nil {
		return cfg, false
	}
	return cfg, ok
}

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
