package web

import (
	"advent2024/pkg/solver"
	"advent2024/web/config"
	"advent2024/web/middleware"
	"advent2024/web/weberrors"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
)

type OAuthReplyGithub struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
type Token interface {
	Token() string
}

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
	logger := middleware.GetLogger(r)

	type registeredKeys struct {
		Keys []string
	}

	b, err := json.Marshal(registeredKeys{Keys: registered_keys})

	rc := http.StatusInternalServerError
	errMsg := "unable to marchal registered solvers"
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func SolveWithUpload(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	cfg, ok := middleware.GetConfig(r)

	var rc int
	var errMsg string

	rc = http.StatusInternalServerError
	errMsg = "configuration error: solve with upload: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyUploadTemplate)
	template, ok := templateVal.(*template.Template)

	rc = http.StatusInternalServerError
	errMsg = "unable to find upload tempate"
	if weberrors.HandleError(w, logger, weberrors.OkToError(!ok || template == nil), rc, errMsg) != nil {
		return
	}

	// TODO: add oauth for mutliple providers

	var provider config.OAuthProvider

	if cfg.OAuth {
		providerName := "github"

		_, exists := cfg.OAuthProviders[providerName]

		rc = http.StatusBadRequest
		errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)

		if weberrors.HandleError(w, logger, weberrors.OkToError(!exists), rc, errMsg) != nil {
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

	err := template.Execute(w, data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render upload template %v", err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	providerName := r.PathValue("provider")
	_, exists := config.OAuthProviders[providerName]

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)
	if weberrors.HandleError(w, logger, weberrors.OkToError(!exists), rc, errMsg) != nil {
		return
	}

	templateVal := r.Context().Value(middleware.ContextKeyCallbackTemplate)
	template, ok := templateVal.(*template.Template)

	rc = http.StatusInternalServerError
	errMsg = "unable to find callback template"
	if weberrors.HandleError(w, logger, weberrors.OkToError(!ok || template == nil), rc, errMsg) != nil {
		return
	}

	provider := config.OAuthProviders[providerName]

	data := struct {
		Endpoint string
	}{
		Endpoint: provider.TokenEndpoint(),
	}

	err := template.Execute(w, data)

	rc = http.StatusInternalServerError
	errMsg = fmt.Sprintf("unable to render callback template %v", err)
	if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
		return
	}
}

func OAuthHandler(w http.ResponseWriter, r *http.Request) {
	var rc int
	var errMsg string

	logger := middleware.GetLogger(r)
	config, ok := middleware.GetConfig(r)

	rc = http.StatusInternalServerError
	errMsg = "configuration error: unable to get config"
	if weberrors.HandleError(w, logger, weberrors.OkToError(ok), rc, errMsg) != nil {
		return
	}

	providerName := r.PathValue("provider")
	_, exists := config.OAuthProviders[providerName]

	rc = http.StatusBadRequest
	errMsg = fmt.Sprintf("unknown Oauth provider %s", providerName)
	if weberrors.HandleError(w, logger, weberrors.OkToError(!exists), rc, errMsg) != nil {
		return
	}

	query := r.URL.Query()
	provider := config.OAuthProviders[providerName]

	switch provider.Name {
	case "github":
		token, err := exchangeCodeForToken(&provider, query.Get("code"))

		rc = http.StatusInternalServerError
		errMsg = fmt.Sprintf("unable to exchange code for token with %s: %v", provider.Name, err)
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jwtToken, err := generateJWT(provider.Name, token.Token(), config.JWTSecret)

		rc = http.StatusInternalServerError
		errMsg = "uanble to create jwt token"
		if weberrors.HandleError(w, logger, err, rc, errMsg) != nil {
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"access_token": jwtToken})

		return
	}
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

func (rep OAuthReplyGithub) Token() string {
	return rep.AccessToken
}
