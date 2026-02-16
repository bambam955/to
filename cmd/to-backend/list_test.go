package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestListCommand(t *testing.T) {
	t.Run("no aliases", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		_, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--list"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "no aliases registered") {
			t.Errorf("expected 'no aliases registered', got: %s", output)
		}
	})

	t.Run("lists aliases alphabetically", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		dirA := t.TempDir()
		dirB := t.TempDir()
		dirC := t.TempDir()

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("charlie", dirC); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.AddAlias("alpha", dirA); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.AddAlias("bravo", dirB); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--list"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d: %s", len(lines), output)
		}

		if !strings.HasPrefix(lines[0], "alpha") {
			t.Errorf("expected first line to start with 'alpha', got: %s", lines[0])
		}
		if !strings.HasPrefix(lines[1], "bravo") {
			t.Errorf("expected second line to start with 'bravo', got: %s", lines[1])
		}
		if !strings.HasPrefix(lines[2], "charlie") {
			t.Errorf("expected third line to start with 'charlie', got: %s", lines[2])
		}
	})

	t.Run("column alignment", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		dirShort := t.TempDir()
		dirLong := t.TempDir()

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("a", dirShort); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.AddAlias("longname", dirLong); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--list"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d: %s", len(lines), output)
		}

		// Both lines should have paths starting at the same column.
		// "a" is padded to match "longname" length (8), plus 2 spaces.
		if !strings.Contains(lines[0], "a") || !strings.Contains(lines[0], dirShort) {
			t.Errorf("expected first line to contain 'a' and path, got: %s", lines[0])
		}
		if !strings.Contains(lines[1], "longname") || !strings.Contains(lines[1], dirLong) {
			t.Errorf("expected second line to contain 'longname' and path, got: %s", lines[1])
		}

		// Check that paths are aligned (start at same index).
		idx0 := strings.Index(lines[0], "/")
		idx1 := strings.Index(lines[1], "/")
		if idx0 != idx1 {
			t.Errorf("expected paths to be aligned at same column, got %d and %d", idx0, idx1)
		}
	})

	t.Run("initializes database if missing", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--list"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "no aliases registered") {
			t.Errorf("expected 'no aliases registered', got: %s", output)
		}
	})
}
