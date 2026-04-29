package config

import (
	"errors"
	"os"
)

// Config holds the runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	SecretPath string
}

// Load reads configuration from environment variables with optional defaults.
func Load() (*Config, error) {
	addr := getEnv("VAULT_ADDR", "http://127.0.0.1:8200")
	if addr == "" {
		return nil, errors.New("config: VAULT_ADDR must not be empty")
	}

	token := getEnv("VAULT_TOKEN", "")
	if token == "" {
		return nil, errors.New("config: VAULT_TOKEN is required")
	}

	return &Config{
		VaultAddr:  addr,
		VaultToken: token,
		Namespace:  getEnv("VAULTPULL_NAMESPACE", ""),
		OutputFile: getEnv("VAULTPULL_OUTPUT", ".env"),
		SecretPath: getEnv("VAULTPULL_SECRET_PATH", "secret/data/app"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
