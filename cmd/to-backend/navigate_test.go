package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestNavigateCommand(t *testing.T) {
	t.Run("successful navigation", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := filepath.Join(tmpDir, "projects")
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			t.Fatalf("failed to create target dir: %v", err)
		}

		// Init database and add alias
		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("proj", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		t.Setenv("TO_DB", dbPath)

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"proj"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v (stderr: %s)", err, stderr.String())
		}

		output := stdout.String()
		expected := "[to] " + targetDir
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got %q", expected, output)
		}
	})

	t.Run("alias not found", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")

		// Init empty database
		_, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}

		t.Setenv("TO_DB", dbPath)

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"nonexistent"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent alias, got nil")
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "not found") {
			t.Errorf("expected stderr to contain 'not found', got %q", errOutput)
		}
	})

	t.Run("directory no longer exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := filepath.Join(tmpDir, "gone")

		// Create directory, add alias, then remove directory
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			t.Fatalf("failed to create target dir: %v", err)
		}

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("gone", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		// Remove the directory
		if err := os.RemoveAll(targetDir); err != nil {
			t.Fatalf("failed to remove target dir: %v", err)
		}

		t.Setenv("TO_DB", dbPath)

		var stdout, stderr bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetErr(&stderr)
		rootCmd.SetArgs([]string{"gone"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for removed directory, got nil")
		}

		errOutput := stderr.String()
		if !strings.Contains(errOutput, "directory no longer exists") {
			t.Errorf("expected stderr to contain 'directory no longer exists', got %q", errOutput)
		}
		if !strings.Contains(errOutput, targetDir) {
			t.Errorf("expected stderr to contain directory path %q, got %q", targetDir, errOutput)
		}
	})
}
