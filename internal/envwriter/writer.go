package envwriter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FilteredWrite writes secrets as KEY=VALUE lines to w.
// If namespace is non-empty, only keys with that prefix are written.
// Returns (written, filtered, error) where filtered is the count of excluded keys.
func FilteredWrite(w io.Writer, secrets map[string]string, namespace string) (int, int, error) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	written := 0
	filtered := 0

	for _, k := range keys {
		if namespace != "" && !strings.HasPrefix(k, namespace+"_") && k != namespace {
			filtered++
			continue
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, secrets[k]); err != nil {
			return written, filtered, fmt.Errorf("envwriter: write error: %w", err)
		}
		written++
	}

	return written, filtered, nil
}
