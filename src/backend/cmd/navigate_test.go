package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestNavigateCommand(t *testing.T) {
	t.Run("successful navigation", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := filepath.Join(tmpDir, "projects")
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			t.Fatalf("failed to create target dir: %v", err)
		}

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
		var emittedPath string
		stubNavigationControlWriter(t, func(path string) error {
			emittedPath = path
			return nil
		})

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"proj"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if emittedPath != targetDir {
			t.Errorf("expected control channel path %q, got %q", targetDir, emittedPath)
		}

		output := stdout.String()
		if output != "" {
			t.Errorf("expected no stdout output for navigation, got %q", output)
		}
	})

	t.Run("alias not found", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")

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
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := filepath.Join(tmpDir, "gone")

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

	t.Run("no arguments", func(t *testing.T) {
		resetFlags(t)
		rootCmd.SetArgs([]string{})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for no arguments, got nil")
		}

		if !strings.Contains(err.Error(), "usage:") {
			t.Errorf("expected usage hint in error, got %q", err.Error())
		}
	})

	t.Run("returns error when control channel is unavailable", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := filepath.Join(tmpDir, "projects")
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			t.Fatalf("failed to create target dir: %v", err)
		}

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
		stubNavigationControlWriter(t, func(path string) error {
			return fmt.Errorf("bad file descriptor")
		})

		rootCmd.SetArgs([]string{"proj"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error when control channel is unavailable, got nil")
		}
		if !strings.Contains(err.Error(), "navigation protocol channel unavailable") {
			t.Fatalf("expected navigation channel error, got %q", err.Error())
		}
	})
}
