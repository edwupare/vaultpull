package envwriter

import (
	"regexp"
	"strings"
)

// sensitivePatterns holds compiled regexes for keys considered sensitive.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)api_key`),
	regexp.MustCompile(`(?i)private_key`),
	regexp.MustCompile(`(?i)credentials`),
}

const redactedValue = "[REDACTED]"

// IsSensitiveKey returns true if the key matches any known sensitive pattern.
func IsSensitiveKey(key string) bool {
	for _, re := range sensitivePatterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactMap returns a copy of the map with sensitive values replaced.
func RedactMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if IsSensitiveKey(k) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactLine redacts the value portion of a KEY=VALUE line if the key is sensitive.
// Lines that do not match the KEY=VALUE format are returned unchanged.
func RedactLine(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return line
	}
	key := strings.TrimSpace(parts[0])
	if IsSensitiveKey(key) {
		return key + "=" + redactedValue
	}
	return line
}
