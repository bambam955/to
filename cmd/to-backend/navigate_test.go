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

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"proj"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
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

		rootCmd.SetArgs([]string{"nonexistent"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent alias, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected error to contain 'not found', got %q", err.Error())
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

		rootCmd.SetArgs([]string{"gone"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for removed directory, got nil")
		}

		if !strings.Contains(err.Error(), "directory no longer exists") {
			t.Errorf("expected error to contain 'directory no longer exists', got %q", err.Error())
		}
		if !strings.Contains(err.Error(), targetDir) {
			t.Errorf("expected error to contain directory path %q, got %q", targetDir, err.Error())
		}
	})
}
