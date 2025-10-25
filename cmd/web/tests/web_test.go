// Tests for web functionality
package tests

import (
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/webhandlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"
)

func TestOAuthCallback(t *testing.T) {
	// create config
	cfg := config.NewConfig()

	// TODO: generic provider?
	providerCfg := config.GithubProvider{}
	providerCfg.ProviderName = "knownProvider"
	cfg.OAuthProviders[providerCfg.Name()] = providerCfg

	// load template
	callbackTemplate := template.Must(template.ParseFiles("./../templates/pages/callback.tmpl"))

	// setup the router
	mux := http.NewServeMux()
	mux.Handle("GET /callback/{provider}",
		middleware.Chain(
			http.HandlerFunc(webhandlers.OAuthCallback),
			middleware.WithConfig(&cfg),
			middleware.WithTemplate(callbackTemplate)))

	// casses
	cases := []struct {
		name string
		url  string
		want int
	}{
		{"known provider", "/callback/knownProvider", http.StatusOK},
		{"unknown provider", "/callback/unknownProvider", http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// create the request
			req := httptest.NewRequest("GET", c.url, nil)
			w := httptest.NewRecorder()

			// call the router
			mux.ServeHTTP(w, req)

			if w.Code != c.want {
				t.Errorf("got %d, want %d", w.Code, c.want)
			}
		})

	}
}

func TestOAuthHandler(t *testing.T) {

	// Mock OAuth provider, modeled based on GitHub
	// 200 OK - for any request containing valid data
	// Bad Request - for incorrectly encoded data
	// Close connection for non specific code - to simulate network issues
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(``))
			return
		}

		code := r.PostFormValue("code")
		_ = r.PostFormValue("client_id")
		_ = r.PostFormValue("client_secret")

		switch code {
		case "invalidCode":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"bad_verification_code", "error_description": "error description", "error_uri": "http://localhost/error.html"}`))
		case "validCode":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"scope":"read","token_type":"bearer","access_token":"validToken"}`))
		case "simulateWrongURL":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(``))
		default:
			hijacker, _ := w.(http.Hijacker)
			conn, _, _ := hijacker.Hijack()
			conn.Close()
		}
	}))
	defer provider.Close()

	// create config
	cfg := config.NewConfig()

	providerCfg := config.GithubProvider{}
	providerCfg.ProviderName = "github"
	providerCfg.ProviderTokenURL = provider.URL + "/token"
	providerCfg.AppClientId = "dummy"
	providerCfg.AppClientSecret = "dummy"

	cfg.OAuthProviders[providerCfg.Name()] = providerCfg

	// setup the router
	mux := http.NewServeMux()
	mux.Handle("POST /oauth/{provider}/token",
		middleware.Chain(
			http.HandlerFunc(webhandlers.OAuthHandler),
			middleware.WithConfig(&cfg)))

	// cases

	cases := []struct {
		name string
		url  string
		want int
	}{
		{"unknown provider", "/oauth/unknownProvider/token", http.StatusBadRequest},
		{"missing code", "/oauth/github/token", http.StatusBadRequest},
		{"github breaks", "/oauth/github/token?code=breakGithub", http.StatusInternalServerError},
		{"invalid code", "/oauth/github/token?code=invalidCode", http.StatusBadRequest},
		{"valid code", "/oauth/github/token?code=validCode", http.StatusOK},
		{"'wrong' code<>exchange URL", "/oauth/github/token?code=simulateWrongURL", http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// create request
			req := httptest.NewRequest("POST", c.url, nil)
			w := httptest.NewRecorder()

			// call the router
			mux.ServeHTTP(w, req)

			if w.Code != c.want {
				t.Errorf("got %d, want %d", w.Code, c.want)
			}
		})

	}
}

func TestToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		want := true

		secret := "jwtSecret"
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", []byte(secret), time.Duration(5*time.Second))
		parsedToken, _ := middleware.ParseToken(jwtTokenStr, []byte(secret))
		if middleware.TokenValid(parsedToken) != want {
			t.Errorf("token invalid")
		}
	})
	t.Run("invalid token", func(t *testing.T) {
		want := false
		secret := "jwtSecret"
		jwtTokenStr := ""
		parsedToken, _ := middleware.ParseToken(jwtTokenStr, []byte(secret))
		if middleware.TokenValid(parsedToken) != want {
			t.Errorf("token valid")
		}
	})
	t.Run("expired token", func(t *testing.T) {
		want := false

		secret := "jwtSecret"
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", []byte(secret), time.Duration(5*time.Second))

		time.Sleep(6 * time.Second)

		parsedToken, _ := middleware.ParseToken(jwtTokenStr, []byte(secret))
		if middleware.TokenValid(parsedToken) != want {
			t.Errorf("token valid")
		}
	})
	t.Run("wrong secret", func(t *testing.T) {
		want := false

		secret := "jwtSecret"
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", []byte(secret), time.Duration(5*time.Second))
		parsedToken, _ := middleware.ParseToken(jwtTokenStr, []byte(secret+"something"))
		if middleware.TokenValid(parsedToken) != want {
			t.Errorf("token valid")
		}
	})
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function didn't panic")
		}
	}()

	mux := http.NewServeMux()
	f := func(w http.ResponseWriter, r *http.Request) {
		panic("problem")
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.Handle("/", http.HandlerFunc(f))
	mux.ServeHTTP(w, req)
}

func TestRecovery(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("function paniced")
		}
	}()

	mux := http.NewServeMux()
	f := func(w http.ResponseWriter, r *http.Request) {
		panic("problem")
	}

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.Handle("/", middleware.RecoveryMiddleware()(http.HandlerFunc(f)))
	mux.ServeHTTP(w, req)
}
