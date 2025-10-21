package tests

import (
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/web"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"
)

func TestOAuthCallback(t *testing.T) {
	// create config
	cfg := config.NewConfig()

	providerCfg := config.OAuthProvider{}
	providerCfg.Name = "knownProvider"
	cfg.OAuthProviders[providerCfg.Name] = providerCfg

	// load template
	callbackTemplate := template.Must(template.ParseFiles("./../templates/callback.tmpl"))

	// setup the router
	mux := http.NewServeMux()
	mux.Handle("GET /callback/{provider}",
		middleware.Chain(
			http.HandlerFunc(web.OAuthCallback),
			middleware.WithConfig(&cfg),
			middleware.WithTemplate(callbackTemplate, middleware.ContextKeyCallbackTemplate)))

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

	// Mock OAuth provider
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(``))
			return
		}

		code := r.PostFormValue("code")
		_ = r.PostFormValue("client_id")
		_ = r.PostFormValue("client_secret")

		switch code {
		case "":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(``))
		case "validCode":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"scope":"read","token_type":"bearer","access_token":"validToken"}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{}`))
		}
	}))
	defer provider.Close()

	// create config
	cfg := config.NewConfig()

	providerCfg := config.OAuthProvider{}
	providerCfg.Name = "github"
	providerCfg.TokenURL = provider.URL + "/token"
	providerCfg.ClientId = "dummy"
	providerCfg.ClientSecret = "dummy"
	// providerCfg.CallbackURL = *oAuthGithubCallbackURL
	// providerCfg.UserAuthURL = *oAuthGithubUserAuthURL

	cfg.OAuthProviders[providerCfg.Name] = providerCfg

	// setup the router
	mux := http.NewServeMux()
	mux.Handle("POST /oauth/{provider}/token",
		middleware.Chain(
			http.HandlerFunc(web.OAuthHandler),
			middleware.WithConfig(&cfg)))

	// cases

	cases := []struct {
		name string
		url  string
		want int
	}{
		{"unknown provider", "/oauth/unknownProvider/token", http.StatusBadRequest},
		{"missing code", "/oauth/github/token", http.StatusBadRequest},
		// TODO: what github returns in case of bad token?
		// 200 OK and you have to read the message
		{"invalid code", "/oauth/github/token?code=invalidCode", http.StatusInternalServerError},
		{"valid code", "/oauth/github/token?code=validCode", http.StatusOK},
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
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", secret, time.Duration(5*time.Second))
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
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", secret, time.Duration(5*time.Second))

		time.Sleep(6 * time.Second)

		parsedToken, _ := middleware.ParseToken(jwtTokenStr, []byte(secret))
		if middleware.TokenValid(parsedToken) != want {
			t.Errorf("token valid")
		}
	})
	t.Run("wrong secret", func(t *testing.T) {
		want := false

		secret := "jwtSecret"
		jwtTokenStr, _ := middleware.GenerateJWT("provider", "token", secret, time.Duration(5*time.Second))
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
