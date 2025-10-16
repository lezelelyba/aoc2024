package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      int
	EnableTLS bool
	CertFile  string
	KeyFile   string
	APIRate   int
	APIBurst  int
}

func NewConfig() Config {
	return Config{
		Port:      8080,
		EnableTLS: false,
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

	return config, errs

}

func envOrDefault(env, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}
