package envwriter

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds details about an invalid env key or value.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid env entry %q: %s", e.Key, e.Message)
}

// ValidationResult contains all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Summary() string {
	if !r.HasErrors() {
		return "all entries valid"
	}
	lines := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		lines = append(lines, e.Error())
	}
	return strings.Join(lines, "\n")
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// ValidateEnv checks that all keys in the provided map conform to
// POSIX env variable naming conventions (uppercase, digits, underscores,
// must start with a letter). Values may not contain raw newlines.
func ValidateEnv(secrets map[string]string) ValidationResult {
	result := ValidationResult{}
	for k, v := range secrets {
		if k == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "key must not be empty",
			})
			continue
		}
		if !validKeyRe.MatchString(k) {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "key must match [A-Z][A-Z0-9_]*",
			})
		}
		if strings.ContainsAny(v, "\n\r") {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "value must not contain newline characters",
			})
		}
	}
	return result
}
