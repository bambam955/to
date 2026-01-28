package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetDatabasePath(t *testing.T) {
	t.Run("with TO_DB override", func(t *testing.T) {
		os.Setenv("TO_DB", "/tmp/custom.db")
		defer os.Unsetenv("TO_DB")

		path, err := GetDatabasePath()
		if err != nil {
			t.Fatalf("GetDatabasePath() returned error: %v", err)
		}
		if path != "/tmp/custom.db" {
			t.Errorf("GetDatabasePath() = %q, want /tmp/custom.db", path)
		}
	})

	t.Run("without TO_DB, uses XDG", func(t *testing.T) {
		os.Unsetenv("TO_DB")

		path, err := GetDatabasePath()
		if err != nil {
			t.Fatalf("GetDatabasePath() returned error: %v", err)
		}

		// Should end with .config/to/database.json
		expected := filepath.Join(".config", "to", "database.json")
		if !strings.HasSuffix(path, expected) {
			t.Errorf("GetDatabasePath() = %q, should end with %q", path, expected)
		}
	})
}

func TestGetConfigDir(t *testing.T) {
	// Use a temporary home directory
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	path, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() returned error: %v", err)
	}

	expectedPath := filepath.Join(tmpHome, ".config", "to")
	if path != expectedPath {
		t.Errorf("GetConfigDir() = %q, want %q", path, expectedPath)
	}

	// Verify directory was created
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Config directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Path exists but is not a directory")
	}

	// Verify it's idempotent
	path2, err := GetConfigDir()
	if err != nil {
		t.Errorf("Second call to GetConfigDir() returned error: %v", err)
	}
	if path != path2 {
		t.Errorf("Idempotency broken: %q != %q", path, path2)
	}
}

func TestEnsureDatabaseDir(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	err := EnsureDatabaseDir()
	if err != nil {
		t.Fatalf("EnsureDatabaseDir() returned error: %v", err)
	}

	// Verify directory was created
	expectedPath := filepath.Join(tmpHome, ".config", "to")
	info, err := os.Stat(expectedPath)
	if err != nil {
		t.Fatalf("Database directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Path exists but is not a directory")
	}
}
