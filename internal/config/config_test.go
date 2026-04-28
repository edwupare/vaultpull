package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestLoad_Success(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "http://vault.example.com:8200")
	setEnv(t, "VAULT_TOKEN", "s.testtoken")
	setEnv(t, "VAULT_NAMESPACE", "myteam")
	setEnv(t, "VAULTPULL_OUTPUT", "secrets.env")
	setEnv(t, "VAULTPULL_PREFIX", "APP_")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("unexpected VaultAddr: %s", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.testtoken" {
		t.Errorf("unexpected VaultToken: %s", cfg.VaultToken)
	}
	if cfg.Namespace != "myteam" {
		t.Errorf("unexpected Namespace: %s", cfg.Namespace)
	}
	if cfg.OutputFile != "secrets.env" {
		t.Errorf("unexpected OutputFile: %s", cfg.OutputFile)
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("unexpected Prefix: %s", cfg.Prefix)
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, "VAULT_TOKEN", "s.token")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_NAMESPACE")
	os.Unsetenv("VAULTPULL_OUTPUT")
	os.Unsetenv("VAULTPULL_PREFIX")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default VaultAddr, got: %s", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile, got: %s", cfg.OutputFile)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:8200")
	os.Unsetenv("VAULT_TOKEN")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN, got nil")
	}
}

func TestLoad_EmptyAddr(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "")
	setEnv(t, "VAULT_TOKEN", "s.token")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for empty VAULT_ADDR, got nil")
	}
}
