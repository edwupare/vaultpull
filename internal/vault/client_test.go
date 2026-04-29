package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, path string, body map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestNewClient_InvalidAddr(t *testing.T) {
	_, err := NewClient("://bad-url", "token", "secret")
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}

func TestReadSecrets_KVv2(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"API_KEY": "abc123",
				"DB_PASS": "secret",
			},
		},
	}
	ts := newTestServer(t, "/v1/secret/data/myapp", payload)
	defer ts.Close()

	c, err := NewClient(ts.URL, "test-token", "secret/data")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	secrets, err := c.ReadSecrets("myapp")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}

	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", secrets["DB_PASS"])
	}
}

func TestReadSecrets_EmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	}))
	defer ts.Close()

	c, err := NewClient(ts.URL, "test-token", "secret")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = c.ReadSecrets("missing")
	if err == nil {
		t.Fatal("expected error for empty secret, got nil")
	}
}

func TestFlattenData(t *testing.T) {
	raw := map[string]interface{}{"FOO": "bar", "NUM": 42}
	out := flattenData(raw)
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
	if out["NUM"] != "42" {
		t.Errorf("expected NUM=42, got %q", out["NUM"])
	}
}
