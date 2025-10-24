package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config
type Config struct {
	Version          string
	Port             int
	EnableTLS        bool
	CertFile         string
	KeyFile          string
	APIRate          int
	APIBurst         int
	SolverTimeout    time.Duration
	OAuth            bool
	JWTSecret        string
	JWTTokenValidity time.Duration
	OAuthProviders   map[string]OAuthProvider
}

// OAuth Provider
type OAuthProvider interface {
	Name() string
	AuthURL() string
	TokenURL() string
	AppCallbackURL() string
	AppTokenEndpoint() string
	ClientID() string
	ClientSecret() string
}

// Github OAuth provider
type GithubProvider struct {
	ProviderName     string
	ProviderAuthURL  string
	ProviderTokenURL string
	CallbackURL      string
	AppClientId      string
	AppClientSecret  string
}

// Constructor with defaults
func NewConfig() Config {
	return Config{
		Port:             8080,
		EnableTLS:        false,
		OAuth:            false,
		APIRate:          3,
		APIBurst:         3,
		SolverTimeout:    time.Duration(5 * time.Second),
		JWTTokenValidity: time.Duration(900 * time.Second),
		OAuthProviders:   make(map[string]OAuthProvider),
	}
}

// Returns Name
func (p GithubProvider) Name() string {
	return p.ProviderName
}

// Returns code<>token exchange URL where client can exchange code for JWT token
func (p GithubProvider) AppTokenEndpoint() string {

	return fmt.Sprintf("/oauth/%s/token", p.ProviderName)
}

// Returns callback URL where Github will redirect the user
func (p GithubProvider) AppCallbackURL() string {
	return p.CallbackURL
}

// Returns URL where client can authenticate with the OAuth Provider
func (p GithubProvider) AuthURL() string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=read:user", p.ProviderAuthURL, p.AppClientId, url.QueryEscape(p.CallbackURL))
}

// Returns URl where app can exchange client code for a token
func (p GithubProvider) TokenURL() string {
	return p.ProviderTokenURL
}

// Returns Client ID
func (p GithubProvider) ClientID() string {
	return p.AppClientId
}

// Returns Client Secret
func (p GithubProvider) ClientSecret() string {
	return p.AppClientSecret
}

// Load configuration from command line and environment
// Returns Configuration and list of load errors
// TODO: JSON config
func LoadConfig() (Config, []error) {
	var errs []error

	config := NewConfig()

	// all params are parsed as strings to be able to utilize the ENV variables as defaults

	port := flag.String("port", envOrDefault("PORT", strconv.Itoa(config.Port)), "TCP port on which the app runs, default 8080")
	// port := flag.Int("port", config.Port, "TCP port on which the app runs")

	enableHttps := flag.String("https", envOrDefault("ENABLE_HTTPS", fmt.Sprintf("%t", config.EnableTLS)), "Enables HTTPS, requires cert and key to be specified")
	// enableHttps := flag.Bool("https", false, "Enables HTTPS, requires cert and key to be specified")
	cert := flag.String("cert", envOrDefault("TLS_CERT_FILE", ""), "cert file")
	key := flag.String("key", envOrDefault("TLS_KEY_FILE", ""), "key file")

	apiRate := flag.String("apirate", envOrDefault("API_RATE", strconv.Itoa(config.APIRate)), "API rate limit per second, default 3")
	// apiRate := flag.Int("apirate", config.APIRate, "API rate limit per second")
	apiBurst := flag.String("apiburst", envOrDefault("API_BURST", strconv.Itoa(config.APIBurst)), "API rate burst size, default 3")
	// apiBurst := flag.Int("apiburst", config.APIBurst, "API rate burst size")
	defVal := int(config.SolverTimeout.Seconds())
	solverTimeout := flag.String("solver-timeout", envOrDefault("API_SOLVER_TIMEOUT", strconv.Itoa(defVal)), "Solver timeout in seconds")
	// solverTimeout := flag.Int("solver-timeout", int(config.SolverTimeout.Seconds()), "Solver timeout in seconds")

	oAuth := flag.String("oauth", envOrDefault("ENABLE_OAUTH", fmt.Sprintf("%t", config.OAuth)), "Enables OAuth API authentication, requireds jwt secret and then requireds per provider callback url, token exhcnage url, client id and secret to be specified")
	jwtSecret := flag.String("jwt-secret", envOrDefault("JWT_SECRET", ""), "JWT Secret")
	defVal = int(config.JWTTokenValidity.Seconds())
	jwtTokenValidity := flag.String("jwt-token-validity", envOrDefault("JWT_TOKEN_VALIDITY", strconv.Itoa(defVal)), "JWT Token Validity in seconds")
	// jwtTokenValidity := flag.String("jwt-token-validity", int(config.JWTTokenValidity.Seconds())), "JWT Token Validity in seconds")
	oAuthGithubCallbackURL := flag.String("oauth-github-callback-url", envOrDefault("OAUTH_GITHUB_CALLBACK_URL", ""), "Github OAuth callback")
	oAuthGithubUserAuthURL := flag.String("oauth-github-user-auth-url", envOrDefault("OAUTH_GITHUB_USER_AUTH_URL", ""), "Github OAuth User auth URL")
	oAuthGithubTokenURL := flag.String("oauth-github-token-url", envOrDefault("OAUTH_GITHUB_TOKEN_URL", ""), "Github OAuth Token exchange URL")
	oAuthGithubId := flag.String("oauth-github-id", envOrDefault("OAUTH_GITHUB_CLIENT_ID", ""), "Github OAuth Client ID")
	oAuthGithubSecret := flag.String("oauth-github-secret", envOrDefault("OAUTH_GITHUB_CLIENT_SECRET", ""), "Github OAuth Secret ID")

	flag.Parse()

	// parse int helper
	parseInt := func(name, value string, dest *int, fallback int) {
		v, err := strconv.Atoi(value)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to parse %s, using default %d", name, fallback))
			*dest = fallback
			return
		}

		*dest = v
	}

	// parse int values
	parseInt("port", *port, &config.Port, config.Port)
	parseInt("apiRate", *apiRate, &config.APIRate, config.APIRate)
	parseInt("apiBurst", *apiBurst, &config.APIBurst, config.APIBurst)

	// parse durations
	var durationInt int
	parseInt("jwtTokenValidity", *jwtTokenValidity, &durationInt, int(config.JWTTokenValidity.Seconds()))
	config.JWTTokenValidity = time.Duration(time.Duration(durationInt) * time.Second)

	parseInt("solverTimeout", *solverTimeout, &durationInt, int(config.SolverTimeout.Seconds()))
	config.SolverTimeout = time.Duration(time.Duration(durationInt) * time.Second)

	// parse https
	if *enableHttps == "true" {
		config.EnableTLS = true
	}

	config.CertFile = *cert
	config.KeyFile = *key

	// parse oauth
	if *oAuth == "true" {
		config.OAuth = true
	}

	config.JWTSecret = *jwtSecret

	// parse individual providers
	// TODO: JSON file
	provider := GithubProvider{ProviderName: "github"}

	provider.CallbackURL = *oAuthGithubCallbackURL
	provider.AppClientId = *oAuthGithubId
	provider.AppClientSecret = *oAuthGithubSecret
	provider.ProviderAuthURL = *oAuthGithubUserAuthURL
	provider.ProviderTokenURL = *oAuthGithubTokenURL

	config.OAuthProviders[provider.Name()] = provider

	return config, errs
}

// Validates configuration.
// Returns if configuration is valid or list of validation errors
func (cfg *Config) ValidateConfig() (bool, []error) {

	var errs []error

	var valid bool = true

	// port range
	if cfg.Port < 0 || cfg.Port > 65535 {
		valid = false
		errs = append(errs, fmt.Errorf("port %d outside of range 0 - 65535", cfg.Port))
	}

	// if TLS is enabled both cert and key has to be provided
	if cfg.EnableTLS {
		// validation breaking errors
		tlsErrors := 0

		checks := []struct {
			path string
			name string
		}{
			{cfg.CertFile, "certificate"},
			{cfg.KeyFile, "key"},
		}

		for _, c := range checks {
			if c.path == "" {
				errs = append(errs, fmt.Errorf("TLS enabled, but %s path is missing", c.name))
				tlsErrors++
				continue
			}

			if _, err := os.Stat(c.path); errors.Is(err, os.ErrNotExist) {
				errs = append(errs, fmt.Errorf("TLS enabled, but %s file not found", c.name))
				tlsErrors++
			}
		}

		// TODO: Validate the key and cert are valid and that they belong together

		if tlsErrors > 0 {
			valid = false
		}
	}

	// if OAuth is enabled, providers have to be valid and JWT token secret has to be set
	if cfg.OAuth {
		// validation breaking errors
		oauthErrors := 0

		for _, provider := range cfg.OAuthProviders {
			// check if the values are not empty
			checks := []struct {
				str  string
				name string
			}{
				{cfg.JWTSecret, "JWT secret"},
				{provider.AuthURL(), "user authentication URL"},
				{provider.TokenURL(), "token exchange URL"},
				{provider.ClientID(), "client id"},
				{provider.ClientSecret(), "client secret"},
				{provider.AppCallbackURL(), "callback URL"},
			}

			for _, c := range checks {
				if c.str == "" {
					errs = append(errs, fmt.Errorf("OAuth enable, but %s is missing", c.name))
					oauthErrors++
				}
			}

			// check if values are valid URLs
			checks = []struct {
				str  string
				name string
			}{
				{provider.AuthURL(), "user authentication URL"},
				{provider.TokenURL(), "token exchange URL"},
				{provider.AppCallbackURL(), "callback URL"},
			}

			for _, c := range checks {

				if !isValidURL(c.str) {
					errs = append(errs, fmt.Errorf("%s: invalid URL", c.name))
					oauthErrors++
				}
			}

			if oauthErrors > 0 {
				valid = false
			}
		}
	}

	return valid, errs
}

// Returns value of env variable or default if not set
func envOrDefault(env, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}

// Check if string is valid URL
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func StripHost(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	// Build only the path+query+fragment part
	out := u.Path
	if u.RawQuery != "" {
		out += "?" + u.RawQuery
	}
	if u.Fragment != "" {
		out += "#" + u.Fragment
	}
	return out, nil
}
