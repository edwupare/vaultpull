package envwriter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added   int
	Updated int
	Unchanged int
}

// MergeEnvFile merges newVars into an existing .env file, preserving
// comments and unrelated keys. Keys present in newVars are added or
// updated; all other lines are left intact.
func MergeEnvFile(path string, newVars map[string]string) (MergeResult, error) {
	var result MergeResult
	existing := map[string]string{}
	var lines []string

	f, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return result, fmt.Errorf("open %s: %w", path, err)
	}
	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
			if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existing[parts[0]] = parts[1]
			}
		}
		f.Close()
		if err := scanner.Err(); err != nil {
			return result, fmt.Errorf("scan %s: %w", path, err)
		}
	}

	updated := map[string]bool{}
	for i, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if val, ok := newVars[key]; ok {
			if val != existing[key] {
				lines[i] = key + "=" + val
				result.Updated++
			} else {
				result.Unchanged++
			}
			updated[key] = true
		}
	}

	for key, val := range newVars {
		if !updated[key] {
			lines = append(lines, key+"="+val)
			result.Added++
		}
	}

	out, err := os.Create(path)
	if err != nil {
		return result, fmt.Errorf("create %s: %w", path, err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return result, w.Flush()
}
