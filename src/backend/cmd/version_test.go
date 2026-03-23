package main

import (
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	t.Run("development build prints dev version", func(t *testing.T) {
		resetFlags(t)
		setBuildVersion(t, "")
		forceTTYDetection(t, true)

		var stdout testTTYBuffer
		stdout.fd = 1
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--version"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "TO dev" {
			t.Fatalf("expected dev version banner, got %q", output)
		}
		if containsANSI(output) {
			t.Fatalf("expected plain version output, got ANSI escapes: %q", output)
		}
	})

	t.Run("release build prints injected version", func(t *testing.T) {
		resetFlags(t)
		setBuildVersion(t, "1.2.3")

		var stdout testTTYBuffer
		stdout.fd = 1
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"-v"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != "TO 1.2.3" {
			t.Fatalf("expected release version banner, got %q", output)
		}
		if containsANSI(output) {
			t.Fatalf("expected plain version output, got ANSI escapes: %q", output)
		}
	})
}
