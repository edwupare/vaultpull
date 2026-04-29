package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func readEnv(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readEnv: %v", err)
	}
	return string(b)
}

func TestMergeEnvFile_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	res, err := MergeEnvFile(path, map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 {
		t.Errorf("want Added=2, got %d", res.Added)
	}
	content := readEnv(t, path)
	if !strings.Contains(content, "FOO=bar") || !strings.Contains(content, "BAZ=qux") {
		t.Errorf("missing keys in output: %s", content)
	}
}

func TestMergeEnvFile_UpdatesExistingKey(t *testing.T) {
	path := writeTempEnv(t, "# comment\nFOO=old\nKEEP=me\n")

	res, err := MergeEnvFile(path, map[string]string{"FOO": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Updated != 1 {
		t.Errorf("want Updated=1, got %d", res.Updated)
	}
	content := readEnv(t, path)
	if !strings.Contains(content, "FOO=new") {
		t.Errorf("key not updated: %s", content)
	}
	if !strings.Contains(content, "KEEP=me") {
		t.Errorf("unrelated key removed: %s", content)
	}
	if !strings.Contains(content, "# comment") {
		t.Errorf("comment removed: %s", content)
	}
}

func TestMergeEnvFile_UnchangedKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=same\n")

	res, err := MergeEnvFile(path, map[string]string{"FOO": "same"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Unchanged != 1 || res.Updated != 0 {
		t.Errorf("want Unchanged=1 Updated=0, got %+v", res)
	}
}

func TestMergeEnvFile_AddsNewKeyToExisting(t *testing.T) {
	path := writeTempEnv(t, "EXISTING=yes\n")

	res, err := MergeEnvFile(path, map[string]string{"NEW_KEY": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 {
		t.Errorf("want Added=1, got %d", res.Added)
	}
	content := readEnv(t, path)
	if !strings.Contains(content, "NEW_KEY=hello") {
		t.Errorf("new key not found: %s", content)
	}
	if !strings.Contains(content, "EXISTING=yes") {
		t.Errorf("existing key removed: %s", content)
	}
}
