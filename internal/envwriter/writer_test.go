package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilteredWrite_NoNamespace(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "xyz"}
	n, err := FilteredWrite(file, secrets, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines written, got %d", n)
	}

	content, _ := os.ReadFile(file)
	if !strings.Contains(string(content), "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", content)
	}
}

func TestFilteredWrite_WithNamespace(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"APP_DB_HOST": "db.local",
		"APP_API_KEY": "secret",
		"OTHER_TOKEN": "ignore",
	}
	n, err := FilteredWrite(file, secrets, "APP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}

	content, _ := os.ReadFile(file)
	if strings.Contains(string(content), "OTHER_TOKEN") {
		t.Errorf("OTHER_TOKEN should be filtered out")
	}
	if !strings.Contains(string(content), "DB_HOST=db.local") {
		t.Errorf("expected stripped key DB_HOST in output")
	}
}

func TestFilteredWrite_EmptyResult(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	secrets := map[string]string{"OTHER_KEY": "value"}
	n, err := FilteredWrite(file, secrets, "APP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}

	// File should not be created when nothing matches
	if _, statErr := os.Stat(file); !os.IsNotExist(statErr) {
		t.Error("expected file to not exist for empty result")
	}
}

func TestFilteredWrite_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	_, err := FilteredWrite(file, secrets, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(file)
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if lines[0] != "A_KEY=a" || lines[1] != "M_KEY=m" || lines[2] != "Z_KEY=z" {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
