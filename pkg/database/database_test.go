package database

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/config"
	"to/pkg/errors"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase()
	if db.Version != config.DefaultDatabaseVersion {
		t.Errorf("Version = %q, want %q", db.Version, config.DefaultDatabaseVersion)
	}
	if db.CreatedAt == "" {
		t.Error("CreatedAt should not be empty")
	}
	if db.UpdatedAt == "" {
		t.Error("UpdatedAt should not be empty")
	}
	if len(db.Aliases) != 0 {
		t.Errorf("Aliases length = %d, want 0", len(db.Aliases))
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "database.json")

	db := NewDatabase()
	testDir := t.TempDir()
	db.Aliases = append(db.Aliases, Alias{
		Name:      "work",
		Directory: testDir,
		CreatedAt: now(),
	})

	if err := db.Save(path); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.Version != db.Version {
		t.Errorf("Version = %q, want %q", loaded.Version, db.Version)
	}
	if len(loaded.Aliases) != 1 {
		t.Fatalf("Aliases length = %d, want 1", len(loaded.Aliases))
	}
	if loaded.Aliases[0].Name != "work" {
		t.Errorf("Alias name = %q, want %q", loaded.Aliases[0].Name, "work")
	}
	if loaded.Aliases[0].Directory != testDir {
		t.Errorf("Alias directory = %q, want %q", loaded.Aliases[0].Directory, testDir)
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/database.json")
	if err == nil {
		t.Fatal("Load() should return error for missing file")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeNotFound {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeNotFound)
	}
}

func TestLoadCorrupted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "database.json")
	os.WriteFile(path, []byte("{invalid json!!!"), 0o644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() should return error for corrupted JSON")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeCorrupted {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeCorrupted)
	}
}

func TestAddAliasValid(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()

	if err := db.AddAlias("work", testDir); err != nil {
		t.Fatalf("AddAlias() returned error: %v", err)
	}
	if len(db.Aliases) != 1 {
		t.Fatalf("Aliases length = %d, want 1", len(db.Aliases))
	}
	if db.Aliases[0].Name != "work" {
		t.Errorf("Alias name = %q, want %q", db.Aliases[0].Name, "work")
	}
	if db.Aliases[0].Directory != testDir {
		t.Errorf("Alias directory = %q, want %q", db.Aliases[0].Directory, testDir)
	}
	if db.Aliases[0].CreatedAt == "" {
		t.Error("Alias CreatedAt should not be empty")
	}
}

func TestAddAliasDuplicate(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()

	db.AddAlias("work", testDir)
	err := db.AddAlias("work", testDir)
	if err == nil {
		t.Fatal("AddAlias() should return error for duplicate name")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeExists {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeExists)
	}
}

func TestAddAliasInvalidName(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()

	tests := []struct {
		name    string
		alias   string
		wantErr bool
	}{
		{"empty", "", true},
		{"starts with hyphen", "-bad", true},
		{"starts with underscore", "_bad", true},
		{"contains space", "my alias", true},
		{"contains dot", "my.alias", true},
		{"contains slash", "my/alias", true},
		{"too long", strings.Repeat("a", 51), true},
		{"valid simple", "work", false},
		{"valid with hyphen", "my-work", false},
		{"valid with underscore", "my_work", false},
		{"valid with numbers", "project123", false},
		{"valid starts with number", "1project", false},
		{"valid single char", "a", false},
		{"max length", strings.Repeat("a", 50), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddAlias(tt.alias, testDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddAlias(%q) error = %v, wantErr %v", tt.alias, err, tt.wantErr)
			}
			if err == nil {
				db.RemoveAlias(tt.alias)
			}
		})
	}
}

func TestAddAliasInvalidDirectory(t *testing.T) {
	db := NewDatabase()

	t.Run("non-absolute path", func(t *testing.T) {
		err := db.AddAlias("test", "relative/path")
		if err == nil {
			t.Fatal("AddAlias() should return error for non-absolute path")
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		err := db.AddAlias("test", "/nonexistent/directory/path")
		if err == nil {
			t.Fatal("AddAlias() should return error for non-existent directory")
		}
	})
}

func TestRemoveAliasExisting(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()
	db.AddAlias("work", testDir)

	if err := db.RemoveAlias("work"); err != nil {
		t.Fatalf("RemoveAlias() returned error: %v", err)
	}
	if len(db.Aliases) != 0 {
		t.Errorf("Aliases length = %d, want 0", len(db.Aliases))
	}
}

func TestRemoveAliasNonExisting(t *testing.T) {
	db := NewDatabase()
	err := db.RemoveAlias("nonexistent")
	if err == nil {
		t.Fatal("RemoveAlias() should return error for non-existing alias")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeNotFound {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeNotFound)
	}
}

func TestGetAliasExisting(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()
	db.AddAlias("work", testDir)

	alias, err := db.GetAlias("work")
	if err != nil {
		t.Fatalf("GetAlias() returned error: %v", err)
	}
	if alias.Name != "work" {
		t.Errorf("Alias name = %q, want %q", alias.Name, "work")
	}
	if alias.Directory != testDir {
		t.Errorf("Alias directory = %q, want %q", alias.Directory, testDir)
	}
}

func TestGetAliasNonExisting(t *testing.T) {
	db := NewDatabase()
	_, err := db.GetAlias("nonexistent")
	if err == nil {
		t.Fatal("GetAlias() should return error for non-existing alias")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeNotFound {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeNotFound)
	}
}

func TestListAliases(t *testing.T) {
	db := NewDatabase()
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	db.AddAlias("work", dir1)
	db.AddAlias("home", dir2)

	aliases := db.ListAliases()
	if len(aliases) != 2 {
		t.Fatalf("ListAliases() length = %d, want 2", len(aliases))
	}

	aliases[0].Name = "modified"
	if db.Aliases[0].Name == "modified" {
		t.Error("ListAliases() should return a copy, not a reference")
	}
}

func TestListAliasesEmpty(t *testing.T) {
	db := NewDatabase()
	aliases := db.ListAliases()
	if len(aliases) != 0 {
		t.Errorf("ListAliases() length = %d, want 0", len(aliases))
	}
}

func TestUpdateLastVisited(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()
	db.AddAlias("work", testDir)

	if err := db.UpdateLastVisited("work"); err != nil {
		t.Fatalf("UpdateLastVisited() returned error: %v", err)
	}

	alias, _ := db.GetAlias("work")
	if alias.LastVisited == "" {
		t.Error("LastVisited should not be empty after update")
	}
}

func TestUpdateLastVisitedNonExisting(t *testing.T) {
	db := NewDatabase()
	err := db.UpdateLastVisited("nonexistent")
	if err == nil {
		t.Fatal("UpdateLastVisited() should return error for non-existing alias")
	}
	toErr, ok := err.(*errors.Error)
	if !ok {
		t.Fatalf("error should be *errors.Error, got %T", err)
	}
	if toErr.Type != errors.ErrorTypeNotFound {
		t.Errorf("error type = %q, want %q", toErr.Type, errors.ErrorTypeNotFound)
	}
}

func TestFindByDirectory(t *testing.T) {
	db := NewDatabase()
	testDir := t.TempDir()
	otherDir := t.TempDir()
	db.AddAlias("alias1", testDir)
	db.AddAlias("alias2", testDir)
	db.AddAlias("alias3", otherDir)

	found := db.FindByDirectory(testDir)
	if len(found) != 2 {
		t.Fatalf("FindByDirectory() length = %d, want 2", len(found))
	}
}

func TestFindByDirectoryNoMatch(t *testing.T) {
	db := NewDatabase()
	found := db.FindByDirectory("/nonexistent")
	if len(found) != 0 {
		t.Errorf("FindByDirectory() length = %d, want 0", len(found))
	}
}

func TestCleanInvalid(t *testing.T) {
	db := NewDatabase()
	validDir := t.TempDir()

	db.Aliases = append(db.Aliases, Alias{
		Name:      "valid",
		Directory: validDir,
		CreatedAt: now(),
	})
	db.Aliases = append(db.Aliases, Alias{
		Name:      "invalid",
		Directory: "/nonexistent/directory/that/does/not/exist",
		CreatedAt: now(),
	})

	removed := db.CleanInvalid()
	if len(removed) != 1 {
		t.Fatalf("CleanInvalid() removed = %d, want 1", len(removed))
	}
	if removed[0].Name != "invalid" {
		t.Errorf("removed alias name = %q, want %q", removed[0].Name, "invalid")
	}
	if len(db.Aliases) != 1 {
		t.Fatalf("Aliases length = %d, want 1", len(db.Aliases))
	}
	if db.Aliases[0].Name != "valid" {
		t.Errorf("remaining alias name = %q, want %q", db.Aliases[0].Name, "valid")
	}
}

func TestCleanInvalidNoneRemoved(t *testing.T) {
	db := NewDatabase()
	validDir := t.TempDir()
	db.Aliases = append(db.Aliases, Alias{
		Name:      "valid",
		Directory: validDir,
		CreatedAt: now(),
	})

	removed := db.CleanInvalid()
	if len(removed) != 0 {
		t.Errorf("CleanInvalid() removed = %d, want 0", len(removed))
	}
}

func TestValidateAliasName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"too long", strings.Repeat("x", 51), true},
		{"starts with hyphen", "-name", true},
		{"starts with underscore", "_name", true},
		{"has spaces", "my name", true},
		{"has dots", "my.name", true},
		{"has special chars", "name@home", true},
		{"valid lowercase", "myalias", false},
		{"valid uppercase", "MyAlias", false},
		{"valid with digits", "alias123", false},
		{"valid digit start", "1alias", false},
		{"valid with hyphens", "my-alias", false},
		{"valid with underscores", "my_alias", false},
		{"valid mixed", "My-Alias_123", false},
		{"valid single char", "a", false},
		{"valid max length", strings.Repeat("a", 50), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAliasName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDirectory(t *testing.T) {
	validDir := t.TempDir()

	t.Run("valid directory", func(t *testing.T) {
		err := ValidateDirectory(validDir)
		if err != nil {
			t.Errorf("ValidateDirectory() returned error: %v", err)
		}
	})

	t.Run("empty path", func(t *testing.T) {
		err := ValidateDirectory("")
		if err == nil {
			t.Error("ValidateDirectory() should return error for empty path")
		}
	})

	t.Run("relative path", func(t *testing.T) {
		err := ValidateDirectory("relative/path")
		if err == nil {
			t.Error("ValidateDirectory() should return error for relative path")
		}
	})

	t.Run("non-existent path", func(t *testing.T) {
		err := ValidateDirectory("/nonexistent/path/abc")
		if err == nil {
			t.Error("ValidateDirectory() should return error for non-existent path")
		}
	})

	t.Run("path too long", func(t *testing.T) {
		longPath := "/" + strings.Repeat("a", 4096)
		err := ValidateDirectory(longPath)
		if err == nil {
			t.Error("ValidateDirectory() should return error for path exceeding max length")
		}
	})

	t.Run("file not directory", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "file.txt")
		os.WriteFile(tmpFile, []byte("test"), 0o644)
		err := ValidateDirectory(tmpFile)
		if err == nil {
			t.Error("ValidateDirectory() should return error for file path")
		}
	})
}

func TestAtomicSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "database.json")

	db := NewDatabase()
	if err := db.Save(path); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("database file not created: %v", err)
	}
	if info.Mode().Perm() != filePermissions {
		t.Errorf("file permissions = %o, want %o", info.Mode().Perm(), filePermissions)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read database file: %v", err)
	}
	var parsed Database
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if parsed.Version != config.DefaultDatabaseVersion {
		t.Errorf("saved version = %q, want %q", parsed.Version, config.DefaultDatabaseVersion)
	}
}

func TestInitDatabase(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "sub", "dir")
	path := filepath.Join(subDir, "database.json")

	db, err := InitDatabase(path)
	if err != nil {
		t.Fatalf("InitDatabase() returned error: %v", err)
	}
	if db.Version != config.DefaultDatabaseVersion {
		t.Errorf("Version = %q, want %q", db.Version, config.DefaultDatabaseVersion)
	}

	if _, err := os.Stat(subDir); err != nil {
		t.Fatalf("directories were not created: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("database file was not created: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() after InitDatabase returned error: %v", err)
	}
	if loaded.Version != config.DefaultDatabaseVersion {
		t.Errorf("loaded version = %q, want %q", loaded.Version, config.DefaultDatabaseVersion)
	}
	if len(loaded.Aliases) != 0 {
		t.Errorf("loaded aliases length = %d, want 0", len(loaded.Aliases))
	}
}

func TestLoadNilAliases(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "database.json")
	data := []byte(`{"version":"1.0","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`)
	os.WriteFile(path, data, 0o644)

	db, err := Load(path)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if db.Aliases == nil {
		t.Error("Aliases should not be nil after loading")
	}
}
