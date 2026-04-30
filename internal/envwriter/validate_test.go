package envwriter

import (
	"strings"
	"testing"
)

func TestValidateEnv_AllValid(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
		"PORT":         "8080",
	}
	result := ValidateEnv(secrets)
	if result.HasErrors() {
		t.Fatalf("expected no errors, got: %s", result.Summary())
	}
}

func TestValidateEnv_LowercaseKey(t *testing.T) {
	secrets := map[string]string{
		"bad_key": "value",
	}
	result := ValidateEnv(secrets)
	if !result.HasErrors() {
		t.Fatal("expected validation error for lowercase key")
	}
	if !strings.Contains(result.Errors[0].Message, "must match") {
		t.Errorf("unexpected error message: %s", result.Errors[0].Message)
	}
}

func TestValidateEnv_KeyStartsWithDigit(t *testing.T) {
	secrets := map[string]string{
		"1INVALID": "value",
	}
	result := ValidateEnv(secrets)
	if !result.HasErrors() {
		t.Fatal("expected validation error for key starting with digit")
	}
}

func TestValidateEnv_EmptyKey(t *testing.T) {
	secrets := map[string]string{
		"": "value",
	}
	result := ValidateEnv(secrets)
	if !result.HasErrors() {
		t.Fatal("expected validation error for empty key")
	}
	if !strings.Contains(result.Errors[0].Message, "must not be empty") {
		t.Errorf("unexpected error message: %s", result.Errors[0].Message)
	}
}

func TestValidateEnv_ValueWithNewline(t *testing.T) {
	secrets := map[string]string{
		"SECRET": "line1\nline2",
	}
	result := ValidateEnv(secrets)
	if !result.HasErrors() {
		t.Fatal("expected validation error for value containing newline")
	}
	if !strings.Contains(result.Errors[0].Message, "newline") {
		t.Errorf("unexpected error message: %s", result.Errors[0].Message)
	}
}

func TestValidateEnv_MultipleErrors(t *testing.T) {
	secrets := map[string]string{
		"bad_key":  "value",
		"ALSO_BAD": "has\nnewline",
	}
	result := ValidateEnv(secrets)
	if len(result.Errors) < 2 {
		t.Fatalf("expected at least 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResult_Summary_NoErrors(t *testing.T) {
	result := ValidationResult{}
	if result.Summary() != "all entries valid" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{Key: "MY_KEY", Message: "some problem"}
	if !strings.Contains(err.Error(), "MY_KEY") {
		t.Errorf("expected key in error string, got: %s", err.Error())
	}
}
