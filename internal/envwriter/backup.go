package envwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupFile creates a timestamped backup of an existing .env file before overwriting it.
// Returns the backup path, or an empty string if the source file did not exist.
func BackupFile(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", path, err)
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}

	ts := time.Now().UTC().Format("20060102T150405Z")
	backupPath := fmt.Sprintf("%s.%s.bak", path, ts)

	if err := os.WriteFile(backupPath, src, 0600); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// CleanupBackups removes backup files older than maxAge for a given .env path.
func CleanupBackups(path string, maxAge time.Duration) error {
	pattern := path + ".*.bak"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob %s: %w", pattern, err)
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(m)
		}
	}
	return nil
}
