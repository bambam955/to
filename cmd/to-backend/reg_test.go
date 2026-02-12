package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegCommand(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "myalias", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "registered 'myalias'") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, targetDir) {
			t.Errorf("expected output to contain directory %s, got: %s", targetDir, output)
		}
	})

	t.Run("duplicate alias", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		// First registration should succeed.
		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "dupalias", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		// Second registration with the same alias name should fail.
		stdout.Reset()
		stderr.Reset()
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "dupalias", targetDir})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for duplicate alias, got nil")
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "already exists") {
			t.Errorf("expected 'already exists' error, got: %s", errOutput)
		}
	})

	t.Run("invalid alias name", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "!!!invalid", targetDir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for invalid alias name, got nil")
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "error:") {
			t.Errorf("expected error message in stderr, got: %s", errOutput)
		}
	})

	t.Run("directory does not exist", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "nodir", "/nonexistent/path"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent directory, got nil")
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "error:") {
			t.Errorf("expected error message in stderr, got: %s", errOutput)
		}
	})

	t.Run("duplicate directory warning", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		// First registration.
		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "first", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		// Second registration with a different alias but the same directory.
		stdout.Reset()
		stderr.Reset()
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "second", targetDir})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("second registration should succeed, got: %v", err)
		}

		// Command should succeed.
		output := stdout.String()
		if !strings.Contains(output, "registered 'second'") {
			t.Errorf("expected success message, got: %s", output)
		}

		// Warning should appear in stderr.
		errOutput := stderr.String()
		if !strings.Contains(errOutput, "warning: directory already registered as: first") {
			t.Errorf("expected duplicate directory warning, got: %s", errOutput)
		}
	})

	t.Run("relative path resolution", func(t *testing.T) {
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		// Create a subdirectory to use as relative path target.
		baseDir := t.TempDir()
		subDir := filepath.Join(baseDir, "subdir")
		if err := os.MkdirAll(subDir, 0o755); err != nil {
			t.Fatalf("failed to create subdir: %v", err)
		}

		// Change working directory so relative path resolves correctly.
		origDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get cwd: %v", err)
		}
		if err := os.Chdir(baseDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		t.Cleanup(func() { os.Chdir(origDir) })

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"reg", "reltest", "subdir"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		// The output should contain the absolute path, not the relative one.
		if !strings.Contains(output, subDir) {
			t.Errorf("expected absolute path %s in output, got: %s", subDir, output)
		}
		if strings.Contains(output, "-> subdir\n") {
			t.Errorf("expected relative path to be resolved, got: %s", output)
		}
	})
}
