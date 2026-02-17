// Package install provides installation utilities for the to tool.
package install

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// BinaryName is the name of the Go backend binary
	BinaryName = "to-backend"
	// WrapperName is the name of the bash wrapper script
	WrapperName = "to.bash"
	// TargetDir is the installation directory
	TargetDir = ".local/bin"
)

// GetInstallPath returns the full installation path for a component.
// It expands ~ to the user's home directory.
func GetInstallPath(component string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	path := filepath.Join(home, TargetDir, component)
	return path, nil
}

// EnsureInstallDir creates the installation directory if it doesn't exist.
// It uses 0755 permissions, which is standard for user-local bin directories.
func EnsureInstallDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	installDir := filepath.Join(home, TargetDir)
	err = os.MkdirAll(installDir, 0o755)
	if err != nil {
		return fmt.Errorf("cannot create install directory %s: %w", installDir, err)
	}

	return nil
}

// VerifyPath checks if the installation directory is in the user's PATH.
// This is primarily informational - the user can fix PATH issues manually.
func VerifyPath() (bool, error) {
	installPath, err := GetInstallPath("")
	if err != nil {
		return false, err
	}

	pathEnv := os.Getenv("PATH")
	return contains(pathEnv, installPath), nil
}

// contains checks if a directory is in a PATH-like string.
func contains(pathEnv, dir string) bool {
	for _, p := range filepath.SplitList(pathEnv) {
		if p == dir {
			return true
		}
	}
	return false
}
