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

type OAuthProvider struct {
	Name         string
	UserAuthURL  string
	TokenURL     string
	CallbackURL  string
	ClientId     string
	ClientSecret string
}

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

func (p OAuthProvider) TokenEndpoint() string {
	return fmt.Sprintf("/oauth/%s/token", p.Name)
}

func (p OAuthProvider) UserAuth() string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=read:user", p.UserAuthURL, p.ClientId, url.QueryEscape(p.CallbackURL))
}

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

	parseInt := func(name, value string, dest *int, fallback int) {
		v, err := strconv.Atoi(value)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to parse %s, using default %d", name, fallback))
			*dest = fallback
			return
		}

		*dest = v
	}

	parseInt("port", *port, &config.Port, config.Port)
	parseInt("apiRate", *apiRate, &config.APIRate, config.APIRate)
	parseInt("apiBurst", *apiBurst, &config.APIBurst, config.APIBurst)

	var durationInt int
	parseInt("jwtTokenValidity", *jwtTokenValidity, &durationInt, int(config.JWTTokenValidity.Seconds()))
	config.JWTTokenValidity = time.Duration(time.Duration(durationInt) * time.Second)

	parseInt("solverTimeout", *solverTimeout, &durationInt, int(config.SolverTimeout.Seconds()))
	config.SolverTimeout = time.Duration(time.Duration(durationInt) * time.Second)

	// replacement for the comment out part
	if *enableHttps == "true" {
		config.EnableTLS = true
	}

	config.CertFile = *cert
	config.KeyFile = *key

	if *oAuth == "true" {
		config.OAuth = true
	}

	config.JWTSecret = *jwtSecret

	provider := OAuthProvider{Name: "github"}

	provider.CallbackURL = *oAuthGithubCallbackURL
	provider.ClientId = *oAuthGithubId
	provider.ClientSecret = *oAuthGithubSecret
	provider.UserAuthURL = *oAuthGithubUserAuthURL
	provider.TokenURL = *oAuthGithubTokenURL

	config.OAuthProviders[provider.Name] = provider

	// if *enableHttps == "true" {

	// 	config.EnableTLS = true
	// 	config.CertFile = *cert
	// 	config.KeyFile = *key

	// 	tlsErrors := 0

	// 	checks := []struct {
	// 		path string
	// 		name string
	// 	}{
	// 		{*cert, "certificate"},
	// 		{*key, "key"},
	// 	}

	// 	for _, c := range checks {
	// 		if c.path == "" {
	// 			errs = append(errs, fmt.Errorf("TLS enabled, but %s path is missing", c.name))
	// 			tlsErrors++
	// 			continue
	// 		}

	// 		if _, err := os.Stat(c.path); errors.Is(err, os.ErrNotExist) {
	// 			errs = append(errs, fmt.Errorf("TLS enabled, but %s file not found", c.name))
	// 			tlsErrors++
	// 		}
	// 	}

	// 	if tlsErrors > 0 {
	// 		config.EnableTLS = false
	// 		errs = append(errs, fmt.Errorf("TLS disabled due to configuration errors"))
	// 	}
	// }

	// if *oAuth == "true" {

	// 	config.OAuth = true
	// 	// config.OAuthProviders = map[string]OAuthProvider{}
	// 	config.JWTSecret = *jwtSecret

	// 	provider := OAuthProvider{Name: "github"}

	// 	provider.CallbackURL = *oAuthGithubCallbackURL
	// 	provider.ClientId = *oAuthGithubId
	// 	provider.ClientSecret = *oAuthGithubSecret
	// 	provider.UserAuthURL = *oAuthGithubUserAuthURL
	// 	provider.TokenURL = *oAuthGithubTokenURL

	// 	oauthErrors := 0

	// 	checks := []struct {
	// 		str  string
	// 		name string
	// 	}{
	// 		{*jwtSecret, "JWT secret"},
	// 		{*oAuthGithubUserAuthURL, "user authentication URL"},
	// 		{*oAuthGithubTokenURL, "token exchange URL"},
	// 		{*oAuthGithubId, "client id"},
	// 		{*oAuthGithubSecret, "client secret"},
	// 		{*oAuthGithubCallbackURL, "callback URL"},
	// 	}

	// 	for _, c := range checks {
	// 		if c.str == "" {
	// 			errs = append(errs, fmt.Errorf("OAuth enable, but %s is missing", c.name))
	// 			oauthErrors++
	// 		}
	// 	}

	// 	if oauthErrors > 0 {
	// 		config.OAuth = false
	// 		errs = append(errs, fmt.Errorf("OAuth disabled due to configuration errors"))
	// 	} else {
	// 		config.OAuthProviders[provider.Name] = provider
	// 	}
	// }

	return config, errs
}

func (cfg *Config) ValidateConfig() (bool, []error) {

	var errs []error

	var valid bool = true

	if cfg.Port < 0 || cfg.Port > 65535 {
		valid = false
		errs = append(errs, fmt.Errorf("port %d outside of range 0 - 65535", cfg.Port))
	}

	if cfg.EnableTLS {
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

		// TODO: Validate the keys and cert are valid and that they belong together

		if tlsErrors > 0 {
			valid = false
		}
	}

	if cfg.OAuth {
		oauthErrors := 0

		for _, provider := range cfg.OAuthProviders {
			checks := []struct {
				str  string
				name string
			}{
				{cfg.JWTSecret, "JWT secret"},
				{provider.UserAuthURL, "user authentication URL"},
				{provider.TokenURL, "token exchange URL"},
				{provider.ClientId, "client id"},
				{provider.ClientSecret, "client secret"},
				{provider.CallbackURL, "callback URL"},
			}

			for _, c := range checks {
				if c.str == "" {
					errs = append(errs, fmt.Errorf("OAuth enable, but %s is missing", c.name))
					oauthErrors++
				}
			}

			checks = []struct {
				str  string
				name string
			}{
				{provider.UserAuthURL, "user authentication URL"},
				{provider.TokenURL, "token exchange URL"},
				{provider.CallbackURL, "callback URL"},
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

func envOrDefault(env, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
