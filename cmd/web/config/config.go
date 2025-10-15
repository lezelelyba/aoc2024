package config

import (
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
	return Config{}
}

func LoadConfig() (Config, []error) {
	var errors []error
	config := NewConfig()

	port := flag.String("port", envOrDefault("PORT", "8080"), "TCP port on which the app runs, default 8080")
	enableHttps := flag.String("https", envOrDefault("ENABLE_HTTPS", "false"), "Enables HTTPS, requires cert and key to be specified")
	cert := flag.String("cert", envOrDefault("TLS_CERT_FILE", ""), "cert file")
	key := flag.String("key", envOrDefault("TLS_KEY_FILE", ""), "key file")
	apiRate := flag.String("apirate", envOrDefault("API_RATE", "3"), "API rate limit per second, default 3")
	apiBurst := flag.String("apiburst", envOrDefault("API_BURST", "3"), "API rate burst size, default 3")

	flag.Parse()

	if port_int, err := strconv.Atoi(*port); err != nil {
		errors = append(errors, err)
		errors = append(errors, fmt.Errorf("unable to parse port, using default %d", config.Port))
	} else {
		config.Port = port_int
	}

	if apiRate_int, err := strconv.Atoi(*apiRate); err != nil {
		errors = append(errors, err)
		errors = append(errors, fmt.Errorf("unable to parse apiRate, using default %d", config.APIRate))
	} else {
		config.APIRate = apiRate_int
	}

	if apiBurst_int, err := strconv.Atoi(*apiBurst); err != nil {
		errors = append(errors, err)
		errors = append(errors, fmt.Errorf("unable to parse apiRate, using default %d", config.APIBurst))
	} else {
		config.APIBurst = apiBurst_int
	}

	if *enableHttps == "true" {
		config.EnableTLS = true

		if *cert == "" {
			errors = append(errors, fmt.Errorf("TLS enabled, but cert path is missing. Disabling TLS"))
			config.EnableTLS = false
		} else {
			config.CertFile = *cert
		}

		if *key == "" {
			errors = append(errors, fmt.Errorf("TLS enabled, but key path is missing. Disabling TLS"))
			config.EnableTLS = false
		} else {
			config.KeyFile = *key
		}
	}

	return config, errors

}

func envOrDefault(env, def string) string {
	if val := os.Getenv(env); val != "" {
		return val
	}
	return def
}
