package tests

import (
	"advent2024/web/config"
	"flag"
	"os"
	"testing"
)

func TestValidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want error
	}{
		{"no args", []string{"app"}, nil},
		{"https", []string{"app",
			"--https", "true",
			"--cert", "test_assets/cert.pem",
			"--key", "test_assets/key.pem",
		}, nil},
		{"OAuth", []string{"app",
			"--oauth", "true",
			"--jwt-secret", "somesecret",
			"--oauth-github-user-auth-url", "https://www.github.com/authorize",
			"--oauth-github-callback-url", "https://localhost/callback",
			"--oauth-github-id", "123",
			"--oauth-github-secret", "123",
			"--oauth-github-token-url", "https://localhost/codeexchange",
		}, nil},
		{"OAuth + https", []string{"app",
			"--https", "true",
			"--cert", "test_assets/cert.pem",
			"--key", "test_assets/key.pem",
			"--oauth", "true",
			"--jwt-secret", "somesecret",
			"--oauth-github-user-auth-url", "https://www.github.com/authorize",
			"--oauth-github-callback-url", "https://localhost/callback",
			"--oauth-github-id", "123",
			"--oauth-github-secret", "123",
			"--oauth-github-token-url", "https://localhost/codeexchange",
		}, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			oldArgs := os.Args
			os.Args = c.args
			defer func() { os.Args = oldArgs }()

			// reset args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, got := config.LoadConfig()
			_, validationGot := cfg.ValidateConfig()

			if len(got) != 0 || len(validationGot) != 0 {
				t.Errorf("want no errors got %d: %v", len(got), got)
			}
		})
	}
}

func TestInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"non numeric port", []string{"app",
			"--port", "asdf",
		}},
		{"non numeric api rate", []string{"app",
			"--apirate", "asdf",
		}},
		{"non numeric api burst", []string{"app",
			"--apiburst", "asdf",
		}},
		{"non numeric validity", []string{"app",
			"--jwt-token-validity", "asdf",
		}},
		{"non numeric solver timeout", []string{"app",
			"--solver-timeout", "asdf",
		}},
		{"https only", []string{"app",
			"--https", "true",
		}},
		{"https unknown files", []string{"app",
			"--https", "true",
			"--cert", "foo.bar",
			"--key", "foo.bar",
		}},
		{"https CA cert", []string{"app",
			"--https", "true",
			"--cert", "test_assets/cert.CA.pem",
			"--key", "test_assets/key.pem",
		}},
		{"https invalid cert", []string{"app",
			"--https", "true",
			"--cert", "test_assets/key.pem",
			"--key", "test_assets/key.pem",
		}},
		{"https invalid key", []string{"app",
			"--https", "true",
			"--cert", "test_assets/cert.pem",
			"--key", "test_assets/cert.pem",
		}},
		{"oauth only", []string{"app",
			"--oauth", "true",
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			oldArgs := os.Args
			os.Args = c.args
			defer func() { os.Args = oldArgs }()

			// reset args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, got := config.LoadConfig()
			_, validationGot := cfg.ValidateConfig()
			if len(got) == 0 && len(validationGot) == 0 {
				t.Errorf("want errors got 0")
			}
		})
	}
}
