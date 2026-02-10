// Package database provides JSON-based storage for directory aliases.
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
	maxAliasNameLength = 50
	maxPathLength      = 4096
	filePermissions    = 0o644
)

var aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

type Alias struct {
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	CreatedAt   string `json:"created_at"`
	LastVisited string `json:"last_visited,omitempty"`
}

type Database struct {
	Version   string  `json:"version"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Aliases   []Alias `json:"aliases"`
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func NewDatabase() *Database {
	ts := now()
	return &Database{
		Version:   config.DefaultDatabaseVersion,
		CreatedAt: ts,
		UpdatedAt: ts,
		Aliases:   []Alias{},
	}
}

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

	if db.Aliases == nil {
		db.Aliases = []Alias{}
	}

	return &db, nil
}

func (db *Database) Save(path string) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to marshal database", err)
	}
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

	if err := tmpFile.Close(); err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, "failed to close temp file", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return errors.Wrap(errors.ErrorTypeDatabase, fmt.Sprintf("failed to save database to: %s", path), err)
	}

	tmpPath = ""
	return nil
}

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

func (db *Database) AddAlias(name, directory string) error {
	if err := ValidateAliasName(name); err != nil {
		return err
	}
	if err := ValidateDirectory(directory); err != nil {
		return err
	}
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

func (db *Database) GetAlias(name string) (*Alias, error) {
	for i := range db.Aliases {
		if db.Aliases[i].Name == name {
			return &db.Aliases[i], nil
		}
	}
	return nil, errors.NotFound(fmt.Sprintf("alias '%s'", name))
}

func (db *Database) ListAliases() []Alias {
	result := make([]Alias, len(db.Aliases))
	copy(result, db.Aliases)
	return result
}

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
		if db.Aliases == nil {
			db.Aliases = []Alias{}
		}
		db.UpdatedAt = now()
	}
	return removed
}

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
