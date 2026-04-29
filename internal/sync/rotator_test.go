package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func newRotatorVaultServer(t *testing.T, secrets map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": secrets,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func newRotatorCfg(addr, token string) *config.Config {
	return &config.Config{
		VaultAddr:  addr,
		VaultToken: token,
		SecretPath: "secret/data/app",
		OutputPath: ".env",
	}
}

func TestRotator_Run_WritesFile(t *testing.T) {
	srv := newRotatorVaultServer(t, map[string]interface{}{"APP_KEY": "secret"})
	defer srv.Close()

	dir := t.TempDir()
	cfg := newRotatorCfg(srv.URL, "test-token")
	cfg.OutputPath = filepath.Join(dir, ".env")

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rotator := NewRotator(s, RotatorConfig{
		VaultPath:  cfg.SecretPath,
		OutputPath: cfg.OutputPath,
		MaxBackups: 3,
	})

	res, err := rotator.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Written != 1 {
		t.Errorf("expected 1 written, got %d", res.Written)
	}
	if _, err := os.Stat(cfg.OutputPath); err != nil {
		t.Errorf("expected output file to exist: %v", err)
	}
}

func TestRotator_Run_DryRun(t *testing.T) {
	srv := newRotatorVaultServer(t, map[string]interface{}{"KEY": "val"})
	defer srv.Close()

	dir := t.TempDir()
	cfg := newRotatorCfg(srv.URL, "test-token")
	cfg.OutputPath = filepath.Join(dir, ".env")

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rotator := NewRotator(s, RotatorConfig{
		VaultPath:  cfg.SecretPath,
		OutputPath: cfg.OutputPath,
		MaxBackups: 3,
		DryRun:     true,
	})

	res, err := rotator.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.BackupPath != "" {
		t.Error("dry run should not produce a backup")
	}
	if _, err := os.Stat(cfg.OutputPath); !os.IsNotExist(err) {
		t.Error("dry run should not write env file")
	}
}

func TestRotator_Run_WithNamespace(t *testing.T) {
	srv := newRotatorVaultServer(t, map[string]interface{}{
		"APP_KEY":   "v1",
		"OTHER_KEY": "v2",
	})
	defer srv.Close()

	dir := t.TempDir()
	cfg := newRotatorCfg(srv.URL, "test-token")
	cfg.OutputPath = filepath.Join(dir, ".env")

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rotator := NewRotator(s, RotatorConfig{
		VaultPath:  cfg.SecretPath,
		OutputPath: cfg.OutputPath,
		MaxBackups: 3,
		Namespace:  "APP",
	})

	res, err := rotator.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Written != 1 {
		t.Errorf("expected 1 written with namespace filter, got %d", res.Written)
	}
}
