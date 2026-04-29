package envwriter

import (
	"strings"
	"testing"
)

func TestDiffEnv_AllAdded(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := DiffEnv(existing, incoming)

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Updated) != 0 || len(result.Removed) != 0 {
		t.Error("expected no updates or removals")
	}
}

func TestDiffEnv_UpdatedKey(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := DiffEnv(existing, incoming)

	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(result.Updated))
	}
	if result.Updated["FOO"] != "new" {
		t.Errorf("expected updated value 'new', got %q", result.Updated["FOO"])
	}
}

func TestDiffEnv_RemovedKey(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD": "gone"}
	incoming := map[string]string{"FOO": "bar"}

	result := DiffEnv(existing, incoming)

	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed))
	}
	if _, ok := result.Removed["OLD"]; !ok {
		t.Error("expected 'OLD' to be in removed")
	}
}

func TestDiffEnv_UnchangedKey(t *testing.T) {
	existing := map[string]string{"FOO": "same"}
	incoming := map[string]string{"FOO": "same"}

	result := DiffEnv(existing, incoming)

	if len(result.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged, got %d", len(result.Unchanged))
	}
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestDiffResult_Summary(t *testing.T) {
	existing := map[string]string{"OLD": "val", "SAME": "x"}
	incoming := map[string]string{"NEW": "val", "SAME": "x", "OLD": "changed"}

	result := DiffEnv(existing, incoming)
	summary := result.Summary()

	if !strings.Contains(summary, "+ NEW") {
		t.Errorf("expected '+ NEW' in summary, got:\n%s", summary)
	}
	if !strings.Contains(summary, "~ OLD") {
		t.Errorf("expected '~ OLD' in summary, got:\n%s", summary)
	}
}

func TestDiffResult_HasChanges_False(t *testing.T) {
	result := &DiffResult{
		Added:     map[string]string{},
		Updated:   map[string]string{},
		Removed:   map[string]string{},
		Unchanged: map[string]string{"A": "1"},
	}
	if result.HasChanges() {
		t.Error("expected HasChanges to return false")
	}
}
