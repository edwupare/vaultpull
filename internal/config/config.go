package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	Prefix     string
}

// Load reads configuration from environment variables and applies defaults.
func Load() (*Config, error) {
	cfg := &Config{
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		Namespace:  getEnv("VAULT_NAMESPACE", ""),
		OutputFile: getEnv("VAULTPULL_OUTPUT", ".env"),
		Prefix:     getEnv("VAULTPULL_PREFIX", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required fields are present.
func (c *Config) validate() error {
	if strings.TrimSpace(c.VaultAddr) == "" {
		return errors.New("VAULT_ADDR must not be empty")
	}
	if strings.TrimSpace(c.VaultToken) == "" {
		return errors.New("VAULT_TOKEN must not be empty")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
