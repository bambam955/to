package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestCleanCommand(t *testing.T) {
	t.Run("no invalid aliases", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		targetDir := t.TempDir()

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("valid", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--clean"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "no invalid aliases found") {
			t.Errorf("expected 'no invalid aliases found', got: %s", output)
		}
	})

	t.Run("removes invalid aliases", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		validDir := t.TempDir()
		goneDir := filepath.Join(tmpDir, "gone")
		if err := os.MkdirAll(goneDir, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("valid", validDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.AddAlias("gone", goneDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		// Remove the directory so it becomes invalid.
		if err := os.RemoveAll(goneDir); err != nil {
			t.Fatalf("failed to remove dir: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--clean"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "removed 'gone'") {
			t.Errorf("expected removal message for 'gone', got: %s", output)
		}
		if !strings.Contains(output, "cleaned 1 alias(es), 1 remaining") {
			t.Errorf("expected summary, got: %s", output)
		}
	})

	t.Run("removes multiple invalid aliases", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		gone1 := filepath.Join(tmpDir, "gone1")
		gone2 := filepath.Join(tmpDir, "gone2")
		if err := os.MkdirAll(gone1, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.MkdirAll(gone2, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("gone1", gone1); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.AddAlias("gone2", gone2); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		if err := os.RemoveAll(gone1); err != nil {
			t.Fatalf("failed to remove dir: %v", err)
		}
		if err := os.RemoveAll(gone2); err != nil {
			t.Fatalf("failed to remove dir: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--clean"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "removed 'gone1'") {
			t.Errorf("expected removal message for 'gone1', got: %s", output)
		}
		if !strings.Contains(output, "removed 'gone2'") {
			t.Errorf("expected removal message for 'gone2', got: %s", output)
		}
		if !strings.Contains(output, "cleaned 2 alias(es), 0 remaining") {
			t.Errorf("expected summary, got: %s", output)
		}
	})

	t.Run("database is persisted after clean", func(t *testing.T) {
		resetFlags(t)
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")
		t.Setenv("TO_DB", dbPath)

		goneDir := filepath.Join(tmpDir, "gone")
		if err := os.MkdirAll(goneDir, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}

		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("gone", goneDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		if err := os.RemoveAll(goneDir); err != nil {
			t.Fatalf("failed to remove dir: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"--clean"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Verify the database was saved with the alias removed.
		reloaded, err := database.Load(dbPath)
		if err != nil {
			t.Fatalf("failed to reload database: %v", err)
		}

		_, err = reloaded.GetAlias("gone")
		if err == nil {
			t.Error("expected alias 'gone' to be removed from database, but it still exists")
		}
	})

	t.Run("empty database", func(t *testing.T) {
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
		rootCmd.SetArgs([]string{"--clean"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "no invalid aliases found") {
			t.Errorf("expected 'no invalid aliases found', got: %s", output)
		}
	})
}
