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
		if !strings.Contains(output, "$ to --list") {
			t.Fatalf("expected help output to include root examples, got: %q", output)
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

func TestHelpExamples(t *testing.T) {
	resetFlags(t)

	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := stdout.String()
	expectedExamples := []struct {
		command     string
		description string
	}{
		{command: "to work", description: `Navigate to the "work" alias`},
		{command: "to --reg work ~/projects/work", description: `Register "work" for a directory`},
		{command: "to --list", description: "List registered aliases"},
		{command: "to --exp work", description: `Print the path for "work"`},
		{command: "to --unreg work", description: `Remove the "work" alias`},
		{command: "to --clean", description: "Remove aliases for missing directories"},
		{command: "to -U", description: "Remove the installed backend and wrappers"},
		{command: "to -P", description: "Remove install artifacts and config data"},
	}
	for _, example := range expectedExamples {
		linePrefix := "  $ " + example.command
		if !helpOutputHasExample(output, linePrefix, example.description) {
			t.Fatalf("expected help output to include example %q with description %q, got: %q", example.command, example.description, output)
		}
	}

	flagsIndex := strings.Index(output, "Flags:")
	examplesIndex := strings.Index(output, "Examples:")
	if flagsIndex == -1 {
		t.Fatalf("expected help output to include Flags section, got: %q", output)
	}
	if examplesIndex == -1 {
		t.Fatalf("expected help output to include Examples section, got: %q", output)
	}
	if examplesIndex < flagsIndex {
		t.Fatalf("expected Examples section after Flags section, got: %q", output)
	}
}

func helpOutputHasExample(output, linePrefix, description string) bool {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, linePrefix) && strings.Contains(line, description) {
			return true
		}
	}
	return false
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
