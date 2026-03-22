package main

import (
	"bytes"
	stderrs "errors"
	"strings"
	"testing"
)

func TestHelpStyling(t *testing.T) {
	t.Run("help output is styled for tty", func(t *testing.T) {
		resetFlags(t)
		forceTTYDetection(t, true)

		originalExample := rootCmd.Example
		rootCmd.Example = "to --list"
		t.Cleanup(func() {
			rootCmd.Example = originalExample
		})

		stdout := &testTTYBuffer{fd: 1}
		rootCmd.SetOut(stdout)
		rootCmd.SetArgs([]string{"--help"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !containsANSI(output) {
			t.Fatalf("expected ANSI escape sequences in TTY help output, got: %q", output)
		}
		if !strings.Contains(output, "Usage:") {
			t.Fatalf("expected help output to include Usage section, got: %q", output)
		}
		if !strings.Contains(output, "Examples:") {
			t.Fatalf("expected help output to include Examples section, got: %q", output)
		}
		if !strings.Contains(output, "Flags:") {
			t.Fatalf("expected help output to include Flags section, got: %q", output)
		}
	})

	t.Run("help output is plain for non-tty", func(t *testing.T) {
		resetFlags(t)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--help"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if containsANSI(output) {
			t.Fatalf("expected non-TTY help output to be plain, got: %q", output)
		}
	})
}

func TestFormatCommandError(t *testing.T) {
	t.Run("tty errors are colored", func(t *testing.T) {
		output := formatCommandError(stderrs.New("boom"), newCLIFormatter(true))
		if !containsANSI(output) {
			t.Fatalf("expected ANSI escape sequences, got: %q", output)
		}
		if !strings.Contains(output, "error: boom") {
			t.Fatalf("expected protocol-compatible error prefix, got: %q", output)
		}
	})

	t.Run("non-tty errors are plain", func(t *testing.T) {
		output := formatCommandError(stderrs.New("boom"), newCLIFormatter(false))
		if output != "error: boom" {
			t.Fatalf("expected plain protocol-compatible error, got: %q", output)
		}
	})
}
