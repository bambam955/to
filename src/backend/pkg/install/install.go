// Package install provides installation utilities for the to tool.
package install

import (
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"

	"to/pkg/config"
	"to/pkg/errors"
)

const (
	// BinaryName is the name of the Go backend binary
	BinaryName = "to-backend"
	// WrapperName is the name of the bash wrapper script
	WrapperName = "to.bash"
	// ZshWrapperName is the name of the zsh wrapper script.
	ZshWrapperName = "to.zsh"
	// FishWrapperName is the name of the fish wrapper script.
	FishWrapperName = "to.fish"
	// TargetDir is the installation directory
	TargetDir = ".local/bin"
)

// ResolveInstallDir returns the directory used for installed binaries and
// wrappers. It prefers the canonical install directory recorded in
// config.toml so uninstall and verify flows stay stable even when the
// original TO_INSTALL_DIR value was relative.
func ResolveInstallDir() (string, error) {
	if recordedDir, err := recordedInstallDir(); err == nil {
		return recordedDir, nil
	} else if !shouldFallbackToEnv(err) {
		return "", err
	}

	if override := os.Getenv("TO_INSTALL_DIR"); override != "" {
		installDir, err := filepath.Abs(override)
		if err != nil {
			return "", fmt.Errorf("cannot resolve install directory %s: %w", override, err)
		}
		return installDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	return filepath.Join(home, TargetDir), nil
}

// recordedInstallDir reads the install directory written by the installer.
// The recorded path is canonical, so backend commands do not have to infer it
// from the current working directory later.
func recordedInstallDir() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	return cfg.InstallDir, nil
}

// shouldFallbackToEnv decides whether a broken install config should be
// ignored so we can still recover the install directory from env/defaults.
// Any config load failure is recoverable because the fallback chain can still
// produce a valid install path.
func shouldFallbackToEnv(err error) bool {
	if err == nil {
		return false
	}
	if stderrors.Is(err, os.ErrNotExist) {
		return true
	}

	var typedErr *errors.Error
	if !stderrors.As(err, &typedErr) {
		return false
	}

	switch typedErr.Type {
	case errors.ErrorTypeInvalid, errors.ErrorTypeCorrupted, errors.ErrorTypePermission, errors.ErrorTypeInternal:
		return true
	default:
		return false
	}
}

// KnownWrapperNames returns the wrapper filenames managed by the installer.
// The uninstall flow removes every known wrapper so switching shells does not
// leave stale entrypoints behind.
func KnownWrapperNames() []string {
	return []string{WrapperName, ZshWrapperName, FishWrapperName}
}

// GetInstallPath returns the full installation path for a component.
// It uses the recorded install directory when available, otherwise falls back
// to TO_INSTALL_DIR or the default user-local bin path.
func GetInstallPath(component string) (string, error) {
	installDir, err := ResolveInstallDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(installDir, component), nil
}

// EnsureInstallDir creates the installation directory if it doesn't exist.
// It uses 0755 permissions, which is standard for user-local bin directories.
func EnsureInstallDir() error {
	installDir, err := ResolveInstallDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return fmt.Errorf("cannot create install directory %s: %w", installDir, err)
	}

	return nil
}

// VerifyPath checks if the installation directory is in the user's PATH.
// This is primarily informational - the user can fix PATH issues manually.
func VerifyPath() (bool, error) {
	installPath, err := ResolveInstallDir()
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
