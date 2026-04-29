package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBackupFile_NoExistingFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")

	bak, err := BackupFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bak != "" {
		t.Errorf("expected empty backup path, got %q", bak)
	}
}

func TestBackupFile_CreatesBackup(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")
	original := []byte("SECRET=abc\nTOKEN=xyz\n")

	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	bak, err := BackupFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bak == "" {
		t.Fatal("expected a backup path, got empty string")
	}
	if !strings.HasSuffix(bak, ".bak") {
		t.Errorf("backup path should end in .bak, got %q", bak)
	}

	got, err := os.ReadFile(bak)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("backup content mismatch: got %q, want %q", got, original)
	}
}

func TestCleanupBackups_RemovesOldFiles(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".env")

	// Create two fake backup files.
	oldBak := path + ".20000101T000000Z.bak"
	newBak := path + ".29991231T235959Z.bak"

	for _, f := range []string{oldBak, newBak} {
		if err := os.WriteFile(f, []byte("x"), 0600); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// Make oldBak appear old by changing its mtime.
	oldTime := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(oldBak, oldTime, oldTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	if err := CleanupBackups(path, 24*time.Hour); err != nil {
		t.Fatalf("CleanupBackups: %v", err)
	}

	if _, err := os.Stat(oldBak); !os.IsNotExist(err) {
		t.Error("expected old backup to be removed")
	}
	if _, err := os.Stat(newBak); err != nil {
		t.Errorf("expected new backup to remain, got: %v", err)
	}
}
