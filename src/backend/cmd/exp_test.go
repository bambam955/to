package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestExpCommand(t *testing.T) {
	t.Run("successful expand", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := t.TempDir()
		t.Setenv("TO_DB", dbPath)

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

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--exp", "proj"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := strings.TrimSpace(stdout.String())
		if output != targetDir {
			t.Errorf("expected %q, got %q", targetDir, output)
		}
		if containsANSI(output) {
			t.Errorf("expected expand output to contain no ANSI escapes, got %q", output)
		}
	})

	t.Run("alias not found", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		_, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}

		rootCmd.SetArgs([]string{"--exp", "nonexistent"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent alias, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected 'not found' in error, got: %s", err.Error())
		}
	})

	t.Run("no arguments", func(t *testing.T) {
		resetFlags(t)

		rootCmd.SetArgs([]string{"--exp"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for no arguments, got nil")
		}

		if !strings.Contains(err.Error(), "usage:") {
			t.Errorf("expected usage hint in error, got: %s", err.Error())
		}
	})

	t.Run("outputs only path", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		targetDir := t.TempDir()
		t.Setenv("TO_DB", dbPath)

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("myproj", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--exp", "myproj"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		// Should be just the path with a newline, nothing else.
		if output != targetDir+"\n" {
			t.Errorf("expected exactly %q, got %q", targetDir+"\n", output)
		}
		if containsANSI(output) {
			t.Errorf("expected expand output to contain no ANSI escapes, got %q", output)
		}
	})
}
