package web

import (
	"advent2024/pkg/solver"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
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
	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	if !ok {
		logger.Printf("solve with upload: unable to get config")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server configuration error"))
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyUploadTemplate)
	template, ok := templateVal.(*template.Template)

	if !ok || template == nil {
		logger.Printf("Unable to find Upload template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	// TODO: add oauth for mutliple providers

	var provider config.OAuthProvider

	if cfg.OAuth {
		providerName := "github"
		if _, exists := cfg.OAuthProviders[providerName]; !exists {
			logger.Printf("unable to find oauth provider %s", providerName)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown OAuth provider"))
			return
		}

		provider = cfg.OAuthProviders[providerName]
	}

	data := struct {
		Title    string
		Endpoint string
		Auth     bool
		Provider *config.OAuthProvider
	}{
		Title:    fmt.Sprintf("Upload Page for day %s", r.PathValue("day")),
		Endpoint: fmt.Sprintf("/api/solvers/%s/%s", r.PathValue("day"), r.PathValue("part")),
		Auth: func() bool {
			if cfg == nil {
				return false
			} else {
				return cfg.OAuth
			}
		}(),
		Provider: &provider,
	}

	if err := template.Execute(w, data); err != nil {
		logger.Printf("Unable to render Upload template %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)

	config, ok := middleware.GetConfig(r)

	if !ok {
		logger.Printf("oauth callback unable to get config")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server configuration error"))
		return
	}

	providerName := r.PathValue("provider")
	if _, exists := config.OAuthProviders[providerName]; !exists {
		logger.Printf("unable to find oauth provider %s", providerName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown OAuth provider"))
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyCallbackTemplate)
	template, ok := templateVal.(*template.Template)

	if !ok || template == nil {
		logger.Printf("Unable to find Callback template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}

	provider := config.OAuthProviders[providerName]

	data := struct {
		Endpoint string
	}{
		Endpoint: provider.TokenEndpoint(),
	}

	if err := template.Execute(w, data); err != nil {
		logger.Printf("Unable to render Callback template")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{}")))
		return
	}
}

func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	if !ok {
		logger.Printf("oauth unable to get config")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server configuration error"))
		return
	}

	providerName := r.PathValue("provider")
	if _, exists := config.OAuthProviders[providerName]; !exists {
		logger.Printf("unable to find oauth provider %s", providerName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown OAuth provider"))
		return
	}

	query := r.URL.Query()
	provider := config.OAuthProviders[providerName]

	switch provider.Name {
	case "github":
		token, err := exchangeCodeForToken(&provider, query.Get("code"))

		if err != nil {
			logger.Printf("unable to resolve token with %s: %v", provider.Name, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unable to exchange code for token with %s", provider.Name)))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jwtToken, err := generateJWT(provider.Name, token.Token(), config.JWTSecret)

		if err != nil {
			logger.Printf("unable to create jwt token")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Unable to create jwt token")))
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"access_token": jwtToken})

		return
	}
}

type OAuthReplyGithub struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (rep OAuthReplyGithub) Token() string {
	return rep.AccessToken
}

type Token interface {
	Token() string
}

func exchangeCodeForToken(provider *config.OAuthProvider, code string) (Token, error) {
	if provider == nil {
		return nil, fmt.Errorf("unable to find empty provider")
	}

	switch provider.Name {
	case "github":
		data := url.Values{}
		data.Set("client_id", provider.ClientId)
		data.Set("client_secret", provider.ClientSecret)
		data.Set("code", code)
		// TODO: redirect_uri

		req, err := http.NewRequest(
			"POST",
			provider.TokenURL,
			strings.NewReader(data.Encode()))

		if err != nil {
			return nil, fmt.Errorf("unable to create OAuth request")
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("unable to exchange code for token with %s", provider.Name)
		}

		var token OAuthReplyGithub

		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal %s response", provider.Name)
		}

		return token, nil

	default:
		return nil, fmt.Errorf("unable to find provider %s", provider.Name)
	}
}

func generateJWT(provider, token, secret string) (string, error) {
	claims := jwt.MapClaims{
		"provider": provider,
		"token":    token,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
