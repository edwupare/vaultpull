package envwriter

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotateEnvFile_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	incoming := map[string]string{
		"APP_KEY":    "abc",
		"APP_SECRET": "xyz",
	}

	res, err := RotateEnvFile(path, incoming, RotateOptions{MaxBackups: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 2 {
		t.Errorf("expected 2 written, got %d", res.Written)
	}
	if len(res.Changes.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Changes.Added))
	}
	if res.BackupPath != "" {
		t.Errorf("expected no backup for new file, got %s", res.BackupPath)
	}
}

func TestRotateEnvFile_UpdatesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("APP_KEY=old\nAPP_SECRET=old\n"), 0600)

	incoming := map[string]string{
		"APP_KEY":    "new",
		"APP_SECRET": "old",
	}

	res, err := RotateEnvFile(path, incoming, RotateOptions{MaxBackups: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(res.Changes.Updated))
	}
	if len(res.Changes.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged, got %d", len(res.Changes.Unchanged))
	}
	if res.BackupPath == "" {
		t.Error("expected backup path to be set")
	}
}

func TestRotateEnvFile_DryRun(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	incoming := map[string]string{"KEY": "val"}
	res, err := RotateEnvFile(path, incoming, RotateOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BackupPath != "" {
		t.Error("dry run should not create backup")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("dry run should not create env file")
	}
}

func TestRotateEnvFile_NamespaceFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	incoming := map[string]string{
		"APP_KEY":   "val1",
		"OTHER_KEY": "val2",
	}

	res, err := RotateEnvFile(path, incoming, RotateOptions{Namespace: "APP", MaxBackups: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 1 {
		t.Errorf("expected 1 written after namespace filter, got %d", res.Written)
	}
}

func TestRotateEnvFile_AuditLogging(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	var buf bytes.Buffer
	logger := &AuditLogger{w: &buf}

	incoming := map[string]string{"DB_URL": "postgres://localhost"}
	_, err := RotateEnvFile(path, incoming, RotateOptions{
		MaxBackups:  3,
		AuditLogger: logger,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_URL") {
		t.Errorf("expected audit log to contain DB_URL, got: %s", buf.String())
	}
}

func TestFilterByNamespace_EmptyNamespace(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	out := filterByNamespace(m, "")
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestFilterByNamespace_WithPrefix(t *testing.T) {
	m := map[string]string{"NS_A": "1", "OTHER_B": "2", "NS_C": "3"}
	out := filterByNamespace(m, "NS")
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}
