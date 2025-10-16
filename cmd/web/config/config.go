package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           int
	EnableTLS      bool
	CertFile       string
	KeyFile        string
	APIRate        int
	APIBurst       int
	OAuth          bool
	OAuthProviders map[string]OAuthProvider
}

type OAuthProvider struct {
	Name         string
	URL          string
	ClientId     string
	ClientSecret string
}

func NewConfig() Config {
	return Config{
		Port:      8080,
		EnableTLS: false,
		OAuth:     false,
		APIRate:   3,
		APIBurst:  3,
	}
}

func LoadConfig() (Config, []error) {
	var errs []error

	config := NewConfig()

	port := flag.String("port", envOrDefault("PORT", strconv.Itoa(config.Port)), "TCP port on which the app runs, default 8080")

	enableHttps := flag.String("https", envOrDefault("ENABLE_HTTPS", fmt.Sprintf("%t", config.EnableTLS)), "Enables HTTPS, requires cert and key to be specified")
	cert := flag.String("cert", envOrDefault("TLS_CERT_FILE", ""), "cert file")
	key := flag.String("key", envOrDefault("TLS_KEY_FILE", ""), "key file")

	apiRate := flag.String("apirate", envOrDefault("API_RATE", strconv.Itoa(config.APIRate)), "API rate limit per second, default 3")
	apiBurst := flag.String("apiburst", envOrDefault("API_BURST", strconv.Itoa(config.APIBurst)), "API rate burst size, default 3")

	oAuth := flag.String("oauth", envOrDefault("ENABLE_OAUTH", fmt.Sprintf("%t", config.OAuth)), "Enables OAuth API authentication, requires client id and secret to be specified")
	oAuthGithubURL := flag.String("oauth-github-url", envOrDefault("OAUTH_GITHUB_CLIENT_URL", ""), "Github OAuth Token exchange URL")
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

	if *enableHttps == "true" {

		config.EnableTLS = true
		config.CertFile = *cert
		config.KeyFile = *key

		tlsErrors := 0

		checks := []struct {
			path string
			name string
		}{
			{*cert, "certificate"},
			{*key, "key"},
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

		if tlsErrors > 0 {
			config.EnableTLS = false
			errs = append(errs, fmt.Errorf("TLS disabled due to configuration errors"))
		}
	}

	if *oAuth == "true" {

		config.OAuth = true
		config.OAuthProviders = map[string]OAuthProvider{}

		oauthErrors := 0

		provider := OAuthProvider{Name: "github"}

		if *oAuthGithubURL == "" {
			errs = append(errs, fmt.Errorf("OAuth enable, but url is missing"))
			oauthErrors++
		} else {
			provider.URL = *oAuthGithubURL
		}

		if *oAuthGithubId == "" {
			errs = append(errs, fmt.Errorf("OAuth enable, but client id is missing"))
			oauthErrors++
		} else {
			provider.ClientId = *oAuthGithubId
		}

		if *oAuthGithubSecret == "" {
			errs = append(errs, fmt.Errorf("OAuth enable, but secret is missing"))
			oauthErrors++
		} else {
			provider.ClientSecret = *oAuthGithubSecret
		}

		if oauthErrors > 0 {
			config.OAuth = false
			errs = append(errs, fmt.Errorf("OAuth disabled due to configuration errors"))
		} else {
			config.OAuthProviders[provider.Name] = provider
		}
	}

	return config, errs

}

func envOrDefault(env, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}
