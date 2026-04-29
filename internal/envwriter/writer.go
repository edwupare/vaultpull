package envwriter

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// FilteredWrite writes secrets to the given file path as KEY=VALUE lines.
// If namespace is non-empty, only keys with that prefix (case-insensitive) are written.
// The prefix is stripped from the key name before writing.
func FilteredWrite(filePath string, secrets map[string]string, namespace string) (int, error) {
	var lines []string
	ns := strings.ToUpper(strings.TrimRight(namespace, "_"))

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		key := strings.ToUpper(k)
		if ns != "" {
			if !strings.HasPrefix(key, ns+"_") {
				continue
			}
			key = strings.TrimPrefix(key, ns+"_")
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, secrets[k]))
	}

	if len(lines) == 0 {
		return 0, nil
	}

	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return 0, fmt.Errorf("envwriter: failed to write %q: %w", filePath, err)
	}

	return len(lines), nil
}
