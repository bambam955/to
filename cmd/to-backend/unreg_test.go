package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/database"
)

func TestUnregCommand(t *testing.T) {
	t.Run("successful unregistration", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")

		t.Setenv("TO_DB", dbPath)

		// Initialize database and add an alias manually.
		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("myalias", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"unreg", "myalias"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "unregistered 'myalias'") {
			t.Errorf("expected success message, got: %q", output)
		}
	})

	t.Run("alias not found", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")

		t.Setenv("TO_DB", dbPath)

		// Initialize an empty database.
		_, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}

		rootCmd.SetArgs([]string{"unreg", "nonexistent"})

		err = rootCmd.Execute()
		if err == nil {
			t.Fatal("expected an error for nonexistent alias, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected 'not found' in error, got: %q", err.Error())
		}
	})

	t.Run("alias is actually removed from database", func(t *testing.T) {
		tmpDir := t.TempDir()
		targetDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "database.json")

		t.Setenv("TO_DB", dbPath)

		// Initialize database and add an alias.
		db, err := database.InitDatabase(dbPath)
		if err != nil {
			t.Fatalf("failed to init database: %v", err)
		}
		if err := db.AddAlias("removeme", targetDir); err != nil {
			t.Fatalf("failed to add alias: %v", err)
		}
		if err := db.Save(dbPath); err != nil {
			t.Fatalf("failed to save database: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"unreg", "removeme"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Reload the database and verify alias is gone.
		reloaded, err := database.Load(dbPath)
		if err != nil {
			t.Fatalf("failed to reload database: %v", err)
		}

		_, err = reloaded.GetAlias("removeme")
		if err == nil {
			t.Error("expected alias 'removeme' to be removed, but it still exists")
		}
	})
}
