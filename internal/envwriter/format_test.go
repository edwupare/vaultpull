package envwriter

import (
	"strings"
	"testing"
)

func TestFormatEnv_Plain(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := FormatEnv(env, FormatOptions{Style: FormatPlain, Sorted: true})
	s := string(out)
	if !strings.Contains(s, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux, got: %s", s)
	}
	if !strings.Contains(s, "FOO=bar\n") {
		t.Errorf("expected FOO=bar, got: %s", s)
	}
}

func TestFormatEnv_Export(t *testing.T) {
	env := map[string]string{"TOKEN": "abc123"}
	out := FormatEnv(env, FormatOptions{Style: FormatExport})
	s := string(out)
	if !strings.Contains(s, "export TOKEN=abc123\n") {
		t.Errorf("expected export prefix, got: %s", s)
	}
}

func TestFormatEnv_Quoted(t *testing.T) {
	env := map[string]string{"MSG": `say "hello"`}
	out := FormatEnv(env, FormatOptions{Style: FormatQuoted})
	s := string(out)
	expected := `MSG="say \"hello\""`
	if !strings.Contains(s, expected) {
		t.Errorf("expected quoted value %q, got: %s", expected, s)
	}
}

func TestFormatEnv_WithComment(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	out := FormatEnv(env, FormatOptions{
		Style:   FormatPlain,
		Comment: "managed by vaultpull",
	})
	s := string(out)
	if !strings.HasPrefix(s, "# managed by vaultpull\n") {
		t.Errorf("expected comment header, got: %s", s)
	}
}

func TestFormatEnv_Sorted(t *testing.T) {
	env := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	out := FormatEnv(env, FormatOptions{Style: FormatPlain, Sorted: true})
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "AAA=") {
		t.Errorf("expected first line AAA=, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZZZ=") {
		t.Errorf("expected last line ZZZ=, got: %s", lines[2])
	}
}

func TestFormatEnv_Empty(t *testing.T) {
	out := FormatEnv(map[string]string{}, FormatOptions{Style: FormatPlain})
	if len(out) != 0 {
		t.Errorf("expected empty output, got: %q", string(out))
	}
}
