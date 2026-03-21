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
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "myalias", targetDir})

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

	t.Run("short flag", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"-r", "shortflag", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "registered 'shortflag'") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("duplicate alias", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "dupalias", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		stdout.Reset()
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "dupalias", targetDir})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for duplicate alias, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("expected 'already exists' error, got: %s", err.Error())
		}
	})

	t.Run("invalid alias name", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		rootCmd.SetArgs([]string{"--reg", "!!!invalid", targetDir})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for invalid alias name, got nil")
		}

		if !strings.Contains(err.Error(), "invalid") {
			t.Errorf("expected 'invalid' in error, got: %s", err.Error())
		}
	})

	t.Run("directory does not exist", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		rootCmd.SetArgs([]string{"--reg", "nodir", "/nonexistent/path"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent directory, got nil")
		}
	})

	t.Run("duplicate directory warning", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "first", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		stdout.Reset()
		var stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"--reg", "second", targetDir})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("second registration should succeed, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "registered 'second'") {
			t.Errorf("expected success message, got: %s", output)
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "warning: directory already registered as: first") {
			t.Errorf("expected duplicate directory warning, got: %s", errOutput)
		}
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		rootCmd.SetArgs([]string{"--reg", "onlyalias"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for wrong arg count, got nil")
		}

		if !strings.Contains(err.Error(), "usage:") {
			t.Errorf("expected usage hint in error, got: %s", err.Error())
		}
	})

	t.Run("relative path resolution", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		baseDir := t.TempDir()
		subDir := filepath.Join(baseDir, "subdir")
		if err := os.MkdirAll(subDir, 0o755); err != nil {
			t.Fatalf("failed to create subdir: %v", err)
		}

		origDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get cwd: %v", err)
		}
		if err := os.Chdir(baseDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}
		t.Cleanup(func() { os.Chdir(origDir) })

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "reltest", "subdir"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, subDir) {
			t.Errorf("expected absolute path %s in output, got: %s", subDir, output)
		}
		if strings.Contains(output, "-> subdir\n") {
			t.Errorf("expected relative path to be resolved, got: %s", output)
		}
	})

	t.Run("alias named reg is allowed", func(t *testing.T) {
		resetFlags(t)
		dbDir := t.TempDir()
		dbPath := filepath.Join(dbDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--reg", "reg", targetDir})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected alias 'reg' to be allowed, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "registered 'reg'") {
			t.Errorf("expected success message, got: %s", output)
		}
	})
}
