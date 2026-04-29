package envwriter

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry represents a single audit log entry for a sync operation.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	Added     int
	Updated   int
	Removed   int
	Unchanged int
}

// String returns a human-readable representation of the audit entry.
func (e AuditEntry) String() string {
	return fmt.Sprintf(
		"[%s] path=%s added=%d updated=%d removed=%d unchanged=%d",
		e.Timestamp.Format(time.RFC3339),
		e.Path,
		e.Added,
		e.Updated,
		e.Removed,
		e.Unchanged,
	)
}

// AuditLogger writes audit entries to a file or writer.
type AuditLogger struct {
	w io.Writer
}

// NewAuditLogger creates an AuditLogger that appends to the given file path.
// If path is empty, logs are written to os.Stdout.
func NewAuditLogger(path string) (*AuditLogger, error) {
	if path == "" {
		return &AuditLogger{w: os.Stdout}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &AuditLogger{w: f}, nil
}

// Log writes an AuditEntry derived from a DiffResult to the logger.
func (a *AuditLogger) Log(envPath string, diff DiffResult) error {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      envPath,
		Added:     len(diff.Added),
		Updated:   len(diff.Updated),
		Removed:   len(diff.Removed),
		Unchanged: len(diff.Unchanged),
	}
	_, err := fmt.Fprintln(a.w, entry.String())
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}
