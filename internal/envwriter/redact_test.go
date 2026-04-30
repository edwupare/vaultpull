package envwriter

import (
	"testing"
)

func TestIsSensitiveKey_Matches(t *testing.T) {
	sensitive := []string{
		"PASSWORD", "db_password", "SECRET", "APP_SECRET",
		"TOKEN", "auth_token", "API_KEY", "PRIVATE_KEY",
		"AWS_CREDENTIALS",
	}
	for _, key := range sensitive {
		if !IsSensitiveKey(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitiveKey_NoMatch(t *testing.T) {
	safe := []string{"HOST", "PORT", "APP_ENV", "LOG_LEVEL", "DATABASE_URL"}
	for _, key := range safe {
		if IsSensitiveKey(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestRedactMap_RedactsSensitiveKeys(t *testing.T) {
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"HOST":        "localhost",
		"PORT":        "5432",
	}
	out := RedactMap(input)

	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT to be unchanged, got %q", out["PORT"])
	}
}

func TestRedactMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"DB_PASSWORD": "supersecret"}
	_ = RedactMap(input)
	if input["DB_PASSWORD"] != "supersecret" {
		t.Error("RedactMap must not mutate the input map")
	}
}

func TestRedactLine_SensitiveLine(t *testing.T) {
	line := "DB_PASSWORD=hunter2"
	got := RedactLine(line)
	expected := "DB_PASSWORD=[REDACTED]"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRedactLine_SafeLine(t *testing.T) {
	line := "HOST=localhost"
	got := RedactLine(line)
	if got != line {
		t.Errorf("expected line unchanged, got %q", got)
	}
}

func TestRedactLine_MalformedLine(t *testing.T) {
	line := "# this is a comment"
	got := RedactLine(line)
	if got != line {
		t.Errorf("expected malformed line unchanged, got %q", got)
	}
}
