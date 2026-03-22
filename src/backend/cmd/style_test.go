package main

import (
	"bytes"
	"testing"
)

func TestIsTTYWriter(t *testing.T) {
	t.Run("writer with fd uses terminal detector", func(t *testing.T) {
		forceTTYDetection(t, true)

		writer := &testTTYBuffer{fd: 1}
		if !isTTYWriter(writer) {
			t.Fatal("expected writer with fd to be detected as TTY")
		}
	})

	t.Run("writer without fd is non-tty", func(t *testing.T) {
		forceTTYDetection(t, true)

		var writer bytes.Buffer
		if isTTYWriter(&writer) {
			t.Fatal("expected buffer without fd to be detected as non-TTY")
		}
	})
}

func TestFormatterEnablement(t *testing.T) {
	t.Run("enabled formatter emits ansi sequences", func(t *testing.T) {
		formatter := newCLIFormatter(true)
		output := formatter.Success("ok")
		if !containsANSI(output) {
			t.Fatalf("expected ANSI style sequence, got %q", output)
		}
	})

	t.Run("disabled formatter is a no-op", func(t *testing.T) {
		formatter := newCLIFormatter(false)
		output := formatter.Success("ok")
		if output != "ok" {
			t.Fatalf("expected plain text output, got %q", output)
		}
	})
}
