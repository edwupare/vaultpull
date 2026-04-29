package envwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateResult holds the outcome of a rotation operation.
type RotateResult struct {
	BackupPath string
	Written    int
	Changes    DiffResult
}

// RotateOptions configures how rotation behaves.
type RotateOptions struct {
	Namespace     string
	MaxBackups    int
	AuditLogger   *AuditLogger
	DryRun        bool
}

// RotateEnvFile performs a full rotate: backup existing file, merge new secrets,
// write the result, and optionally log the diff via the audit logger.
func RotateEnvFile(path string, incoming map[string]string, opts RotateOptions) (*RotateResult, error) {
	maxBackups := opts.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 5
	}

	// Read existing env for diff purposes.
	existing := map[string]string{}
	if data, err := os.ReadFile(path); err == nil {
		for k, v := range parseEnvBytes(data) {
			existing[k] = v
		}
	}

	// Filter incoming by namespace.
	filtered := filterByNamespace(incoming, opts.Namespace)

	// Compute diff before writing.
	diff := DiffEnv(existing, filtered)

	result := &RotateResult{
		Changes: diff,
		Written: len(filtered),
	}

	if opts.DryRun {
		return result, nil
	}

	// Backup existing file.
	backupPath, err := BackupFile(path, time.Now())
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("rotate: backup failed: %w", err)
	}
	result.BackupPath = backupPath

	// Merge and write.
	if err := MergeEnvFile(path, filtered); err != nil {
		return nil, fmt.Errorf("rotate: merge failed: %w", err)
	}

	// Cleanup old backups.
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	if err := CleanupBackups(dir, base, maxBackups); err != nil {
		return nil, fmt.Errorf("rotate: cleanup failed: %w", err)
	}

	// Audit log.
	if opts.AuditLogger != nil {
		for k := range diff.Added {
			opts.AuditLogger.Log("rotate", k, "added")
		}
		for k := range diff.Updated {
			opts.AuditLogger.Log("rotate", k, "updated")
		}
		for k := range diff.Removed {
			opts.AuditLogger.Log("rotate", k, "removed")
		}
	}

	return result, nil
}

// filterByNamespace returns entries from m whose keys have the given namespace prefix.
// If namespace is empty, all entries are returned.
func filterByNamespace(m map[string]string, ns string) map[string]string {
	if ns == "" {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	out := map[string]string{}
	prefix := ns + "_"
	for k, v := range m {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			out[k] = v
		}
	}
	return out
}

// parseEnvBytes parses a simple KEY=VALUE env file byte slice.
func parseEnvBytes(data []byte) map[string]string {
	out := map[string]string{}
	lines := splitLines(string(data))
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, c := range line {
			if c == '=' {
				out[line[:i]] = line[i+1:]
				break
			}
		}
	}
	return out
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
