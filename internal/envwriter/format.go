package envwriter

import (
	"fmt"
	"sort"
	"strings"
)

// FormatStyle defines the output format for env files.
type FormatStyle int

const (
	// FormatExport prefixes each line with "export ".
	FormatExport FormatStyle = iota
	// FormatPlain writes KEY=VALUE without export prefix.
	FormatPlain
	// FormatQuoted wraps values in double quotes.
	FormatQuoted
)

// FormatOptions controls how env entries are rendered.
type FormatOptions struct {
	Style   FormatStyle
	Sorted  bool
	Comment string // optional header comment
}

// FormatEnv renders a map of env vars into a byte slice according to opts.
func FormatEnv(env map[string]string, opts FormatOptions) []byte {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.Sorted {
		sort.Strings(keys)
	}

	var sb strings.Builder

	if opts.Comment != "" {
		for _, line := range strings.Split(opts.Comment, "\n") {
			sb.WriteString("# ")
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		sb.WriteByte('\n')
	}

	for _, k := range keys {
		v := env[k]
		switch opts.Style {
		case FormatExport:
			fmt.Fprintf(&sb, "export %s=%s\n", k, v)
		case FormatQuoted:
			v = strings.ReplaceAll(v, `"`, `\"`)
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		default:
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	return []byte(sb.String())
}
