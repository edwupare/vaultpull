package envwriter

import (
	"fmt"
	"sort"
	"strings"
)

// DiffResult holds the changes between existing and incoming env values.
type DiffResult struct {
	Added   map[string]string
	Updated map[string]string
	Removed map[string]string
	Unchanged map[string]string
}

// HasChanges returns true if there are any additions, updates, or removals.
func (d *DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Updated) > 0 || len(d.Removed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d *DiffResult) Summary() string {
	var sb strings.Builder
	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	for _, k := range keys(d.Added) {
		fmt.Fprintf(&sb, "  + %s\n", k)
	}
	for _, k := range keys(d.Updated) {
		fmt.Fprintf(&sb, "  ~ %s\n", k)
	}
	for _, k := range keys(d.Removed) {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return sb.String()
}

// DiffEnv computes the diff between existing env vars and incoming ones.
func DiffEnv(existing, incoming map[string]string) *DiffResult {
	result := &DiffResult{
		Added:     make(map[string]string),
		Updated:   make(map[string]string),
		Removed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range incoming {
		if oldVal, exists := existing[k]; !exists {
			result.Added[k] = v
		} else if oldVal != v {
			result.Updated[k] = v
		} else {
			result.Unchanged[k] = v
		}
	}

	for k, v := range existing {
		if _, exists := incoming[k]; !exists {
			result.Removed[k] = v
		}
	}

	return result
}
