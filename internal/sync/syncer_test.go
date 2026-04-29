package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vaultpull/internal/config"
	syncer "github.com/vaultpull/internal/sync"
)

func newVaultServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func newCfg(addr, token string) *config.Config {
	return &config.Config{
		VaultAddr:  addr,
		VaultToken: token,
	}
}

func TestRun_WritesEnvFile(t *testing.T) {
	ts := newVaultServer(t, map[string]interface{}{
		"APP_KEY": "value1",
		"DB_HOST": "localhost",
	})
	defer ts.Close()

	cfg := newCfg(ts.URL, "test-token")
	s, err := syncer.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env")
	result, err := s.Run("secret/data/myapp", "", out)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Written != 2 {
		t.Errorf("expected 2 written, got %d", result.Written)
	}

	contents, _ := os.ReadFile(out)
	if !strings.Contains(string(contents), "APP_KEY=value1") {
		t.Errorf("expected APP_KEY in output, got: %s", contents)
	}
}

func TestRun_WithNamespaceFilter(t *testing.T) {
	ts := newVaultServer(t, map[string]interface{}{
		"APP_KEY":  "value1",
		"DB_HOST":  "localhost",
		"APP_PORT": "8080",
	})
	defer ts.Close()

	cfg := newCfg(ts.URL, "test-token")
	s, err := syncer.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env")
	result, err := s.Run("secret/data/myapp", "APP", out)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Written != 2 {
		t.Errorf("expected 2 written after filter, got %d", result.Written)
	}
	if result.Filtered != 1 {
		t.Errorf("expected 1 filtered, got %d", result.Filtered)
	}
}

func TestNew_InvalidAddr(t *testing.T) {
	cfg := newCfg("://bad-addr", "token")
	_, err := syncer.New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid addr, got nil")
	}
}
