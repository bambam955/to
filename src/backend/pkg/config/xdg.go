// Package config provides configuration and XDG Base Directory support.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"to/pkg/errors"
)

const (
	// DatabaseName is the filename of the database
	DatabaseName = "database.json"
	// ConfigDir is the XDG config directory name
	ConfigDir = "to"
	// DefaultDatabaseVersion is the initial schema version
	DefaultDatabaseVersion = "1.0"
)

// GetDatabasePath returns the path to the database file.
// It respects the TO_DB environment variable for overrides.
// Uses XDG Base Directory Specification: ~/.config/to/database.json
func GetDatabasePath() (string, error) {
	// Check for explicit override
	if override := os.Getenv("TO_DB"); override != "" {
		return override, nil
	}

	// Use XDG config directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(errors.ErrorTypeInternal, "cannot determine home directory", err)
	}

	configPath := filepath.Join(home, ".config", ConfigDir)
	return filepath.Join(configPath, DatabaseName), nil
}

// GetConfigDir returns the configuration directory.
// Creates it if it doesn't exist.
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(errors.ErrorTypeInternal, "cannot determine home directory", err)
	}

	configDir := filepath.Join(home, ".config", ConfigDir)
	err = os.MkdirAll(configDir, 0o755)
	if err != nil {
		return "", errors.Wrap(errors.ErrorTypePermission,
			fmt.Sprintf("cannot create config directory %s", configDir), err)
	}

	return configDir, nil
}

// EnsureDatabaseDir ensures the database directory exists.
func EnsureDatabaseDir() error {
	_, err := GetConfigDir()
	return err
}
