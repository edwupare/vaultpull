package envwriter

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditEntry_String(t *testing.T) {
	e := AuditEntry{
		Path:      ".env",
		Added:     2,
		Updated:   1,
		Removed:   0,
		Unchanged: 3,
	}
	s := e.String()
	for _, want := range []string{"path=.env", "added=2", "updated=1", "removed=0", "unchanged=3"} {
		if !strings.Contains(s, want) {
			t.Errorf("AuditEntry.String() missing %q, got: %s", want, s)
		}
	}
}

func TestNewAuditLogger_Stdout(t *testing.T) {
	logger, err := NewAuditLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewAuditLogger_InvalidPath(t *testing.T) {
	_, err := NewAuditLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestAuditLogger_Log_WritesToBuffer(t *testing.T) {
	var buf bytes.Buffer
	logger := &AuditLogger{w: &buf}

	diff := DiffResult{
		Added:     map[string]string{"NEW_KEY": "val"},
		Updated:   map[string][2]string{"OLD_KEY": {"a", "b"}},
		Removed:   map[string]string{},
		Unchanged: map[string]string{"SAME": "x"},
	}

	if err := logger.Log(".env", diff); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	line := buf.String()
	for _, want := range []string{"path=.env", "added=1", "updated=1", "removed=0", "unchanged=1"} {
		if !strings.Contains(line, want) {
			t.Errorf("log line missing %q, got: %s", want, line)
		}
	}
}

func TestNewAuditLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	logger, err := NewAuditLogger(logPath)
	if err != nil {
		t.Fatalf("NewAuditLogger() error: %v", err)
	}

	diff := DiffResult{
		Added:     map[string]string{"A": "1"},
		Updated:   map[string][2]string{},
		Removed:   map[string]string{},
		Unchanged: map[string]string{},
	}
	if err := logger.Log(".env", diff); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if !strings.Contains(string(data), "added=1") {
		t.Errorf("expected 'added=1' in log file, got: %s", string(data))
	}
}
