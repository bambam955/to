// Package database provides JSON-based storage for directory aliases.
//
// The database is loaded entirely into memory for operations and saved
// atomically using a temp-file-and-rename strategy to prevent corruption.
// See spec 002-database-design.md for the full design rationale.
package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"to/pkg/config"
	"to/pkg/errors"
)

const (
	// maxAliasNameLength is the maximum number of characters allowed in an alias name.
	maxAliasNameLength = 50
	// maxPathLength is the maximum number of characters allowed in a directory path.
	maxPathLength = 4096
	// filePermissions is the Unix permission mode applied to the database file.
	filePermissions = 0o644
)

// aliasNamePattern enforces that alias names start with an alphanumeric
// character and may only contain alphanumeric characters, hyphens, or
// underscores.
var aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// Alias represents a single directory alias entry in the database.
type Alias struct {
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	CreatedAt   string `json:"created_at"`
	LastVisited string `json:"last_visited,omitempty"`
}

// Database represents the top-level JSON database structure, containing
// schema version metadata and the collection of aliases.
type Database struct {
	Version   string  `json:"version"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Aliases   []Alias `json:"aliases"`
}

// now returns the current UTC time formatted as an ISO 8601 / RFC 3339 string.
func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// NewDatabase creates a new in-memory database with the default schema version,
// current timestamps, and an empty alias list.
func NewDatabase() *Database {
	ts := now()
	return &Database{
		Version:   config.DefaultDatabaseVersion,
		CreatedAt: ts,
		UpdatedAt: ts,
		Aliases:   []Alias{},
	}
}

// Load reads and parses the database JSON file at the given path.
// Returns categorized errors for missing files, permission issues, and
// corrupted (malformed) JSON.
func Load(path string) (*Database, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(errors.ErrorTypeNotFound, fmt.Sprintf("database file not found: %s", path), err)
		}
		if os.IsPermission(err) {
			return nil, errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot read database file: %s", path), err)
		}
		return nil, errors.Wrap(errors.ErrorTypeDatabase, fmt.Sprintf("failed to read database file: %s", path), err)
	}

	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, errors.Wrap(errors.ErrorTypeCorrupted, fmt.Sprintf("database file is corrupted: %s", path), err)
	}

	// Ensure aliases is never nil so callers can safely iterate without
	// nil checks, even if the JSON had no "aliases" key.
	if db.Aliases == nil {
		db.Aliases = []Alias{}
	}

	return &db, nil
}

// Save writes the database to disk using an atomic write strategy:
//  1. Marshal the database to pretty-printed JSON.
//  2. Write to a temporary file in the same directory.
//  3. Set file permissions.
//  4. Atomically rename the temp file to the final path.
//
// If any step fails, the temp file is cleaned up and the original file
// remains untouched, preventing corruption from partial writes.
func (db *Database) Save(path string) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to marshal database", err)
	}
	// Append a trailing newline for POSIX compatibility.
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "database-*.tmp")
	if err != nil {
		if os.IsPermission(err) {
			return errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot write to directory: %s", dir), err)
		}
		return errors.Wrap(errors.ErrorTypeDatabase, fmt.Sprintf("failed to create temp file in: %s", dir), err)
	}
	tmpPath := tmpFile.Name()

	// Deferred cleanup: remove the temp file if the rename didn't happen.
	// tmpPath is set to "" after a successful rename to skip removal.
	defer func() {
		if tmpPath != "" {
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to write database", err)
	}

	if err := tmpFile.Chmod(filePermissions); err != nil {
		tmpFile.Close()
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to set file permissions", err)
	}

	// Close before rename to flush all data to disk.
	if err := tmpFile.Close(); err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to close temp file", err)
	}

	// Atomic rename replaces the old file in one operation.
	if err := os.Rename(tmpPath, path); err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, fmt.Sprintf("failed to save database to: %s", path), err)
	}

	// Mark rename as successful so the deferred cleanup is skipped.
	tmpPath = ""
	return nil
}

// ValidateAliasName checks that an alias name meets all naming rules:
//   - Must not be empty
//   - Must not exceed maxAliasNameLength characters
//   - Must match aliasNamePattern (starts with alphanumeric, contains only
//     alphanumeric, hyphens, or underscores)
func ValidateAliasName(name string) error {
	if name == "" {
		return errors.InvalidInput("alias name cannot be empty")
	}
	if len(name) > maxAliasNameLength {
		return errors.Newf(errors.ErrorTypeInvalid, "invalid input: alias name exceeds %d characters", maxAliasNameLength)
	}
	if !aliasNamePattern.MatchString(name) {
		return errors.InvalidInput("alias name must start with alphanumeric and contain only alphanumeric, hyphens, or underscores")
	}
	return nil
}

// ValidateDirectory checks that a directory path is valid for alias storage:
//   - Must not be empty
//   - Must not exceed maxPathLength characters
//   - Must be an absolute path
//   - Must exist on disk and be a directory (not a regular file)
//   - Must be accessible (proper permissions)
func ValidateDirectory(directory string) error {
	if directory == "" {
		return errors.InvalidInput("directory path cannot be empty")
	}
	if len(directory) > maxPathLength {
		return errors.Newf(errors.ErrorTypeInvalid, "invalid input: directory path exceeds %d characters", maxPathLength)
	}
	if !filepath.IsAbs(directory) {
		return errors.InvalidInput("directory path must be absolute")
	}
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Newf(errors.ErrorTypeNotFound, "directory does not exist: %s", directory)
		}
		if os.IsPermission(err) {
			return errors.PermissionDenied(directory)
		}
		return errors.Wrap(errors.ErrorTypeDatabase, fmt.Sprintf("cannot access directory: %s", directory), err)
	}
	if !info.IsDir() {
		return errors.Newf(errors.ErrorTypeInvalid, "invalid input: path is not a directory: %s", directory)
	}
	return nil
}

// AddAlias validates the name and directory, checks for duplicates, and
// appends a new alias to the database. The directory path is normalized
// via filepath.Clean before storage for consistent lookups.
func (db *Database) AddAlias(name, directory string) error {
	if err := ValidateAliasName(name); err != nil {
		return err
	}
	if err := ValidateDirectory(directory); err != nil {
		return err
	}
	// Check for duplicate alias names.
	for _, a := range db.Aliases {
		if a.Name == name {
			return errors.AlreadyExists(fmt.Sprintf("alias '%s'", name))
		}
	}
	db.Aliases = append(db.Aliases, Alias{
		Name:      name,
		Directory: filepath.Clean(directory),
		CreatedAt: now(),
	})
	db.UpdatedAt = now()
	return nil
}

// RemoveAlias deletes an alias by name. Returns a NotFound error if the
// alias does not exist.
func (db *Database) RemoveAlias(name string) error {
	for i, a := range db.Aliases {
		if a.Name == name {
			db.Aliases = append(db.Aliases[:i], db.Aliases[i+1:]...)
			db.UpdatedAt = now()
			return nil
		}
	}
	return errors.NotFound(fmt.Sprintf("alias '%s'", name))
}

// GetAlias looks up an alias by name and returns a pointer to it.
// The pointer refers to the actual slice element, so mutations are reflected
// in the database. Returns a NotFound error if the alias does not exist.
func (db *Database) GetAlias(name string) (*Alias, error) {
	for i := range db.Aliases {
		if db.Aliases[i].Name == name {
			return &db.Aliases[i], nil
		}
	}
	return nil, errors.NotFound(fmt.Sprintf("alias '%s'", name))
}

// ListAliases returns a shallow copy of all aliases. The copy ensures
// callers cannot accidentally mutate the database's internal state.
func (db *Database) ListAliases() []Alias {
	result := make([]Alias, len(db.Aliases))
	copy(result, db.Aliases)
	return result
}

// UpdateLastVisited sets the last_visited timestamp of the named alias to
// the current UTC time. Returns a NotFound error if the alias does not exist.
func (db *Database) UpdateLastVisited(name string) error {
	for i := range db.Aliases {
		if db.Aliases[i].Name == name {
			db.Aliases[i].LastVisited = now()
			db.UpdatedAt = now()
			return nil
		}
	}
	return errors.NotFound(fmt.Sprintf("alias '%s'", name))
}

// FindByDirectory returns all aliases that point to the given directory.
// The directory path is normalized before comparison so that equivalent
// paths (e.g., with trailing slashes) match correctly.
func (db *Database) FindByDirectory(directory string) []Alias {
	var result []Alias
	cleaned := filepath.Clean(directory)
	for _, a := range db.Aliases {
		if a.Directory == cleaned {
			result = append(result, a)
		}
	}
	return result
}

// CleanInvalid removes all aliases whose directories no longer exist on
// disk (or are no longer directories). Returns the list of removed aliases
// so the caller can report what was cleaned up.
func (db *Database) CleanInvalid() []Alias {
	var valid []Alias
	var removed []Alias
	for _, a := range db.Aliases {
		if info, err := os.Stat(a.Directory); err != nil || !info.IsDir() {
			removed = append(removed, a)
		} else {
			valid = append(valid, a)
		}
	}
	if len(removed) > 0 {
		db.Aliases = valid
		// Ensure aliases is never nil, even if all entries were removed.
		if db.Aliases == nil {
			db.Aliases = []Alias{}
		}
		db.UpdatedAt = now()
	}
	return removed
}

// InitDatabase creates a new database file at the given path, including
// any missing parent directories (with 0755 permissions). Returns the
// newly created database ready for use.
func InitDatabase(path string) (*Database, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot create directory: %s", dir), err)
	}

	db := NewDatabase()
	if err := db.Save(path); err != nil {
		return nil, err
	}
	return db, nil
}
