package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"to/pkg/config"
	"to/pkg/install"
)

const (
	uninstallUsage   = "usage: to --uninstall"
	purgeUsage       = "usage: to --purge"
	uninstallMessage = "removed TO backend and shell wrappers"
	purgeMessage     = "purged TO database"
)

// runUninstall removes the installed backend binary and all known wrapper
// entrypoints. Missing files are treated as a successful no-op so the command
// can be rerun safely.
func runUninstall(cmd *cobra.Command, args []string) error {
	formatter := formatterForOutput(cmd)

	if len(args) != 0 {
		return fmt.Errorf(uninstallUsage)
	}

	if err := removeInstallArtifacts(); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), formatter.Success(uninstallMessage))
	return nil
}

// runPurge removes the install artifacts and then clears the default config
// directory. Missing files or directories are treated as a successful no-op.
func runPurge(cmd *cobra.Command, args []string) error {
	formatter := formatterForOutput(cmd)

	if len(args) != 0 {
		return fmt.Errorf(purgeUsage)
	}

	if err := removeInstallArtifacts(); err != nil {
		return err
	}

	configDir, err := defaultConfigDir()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("failed to remove config directory %s: %w", configDir, err)
	}

	// Purge is uninstall plus database cleanup, so emit both success messages
	// in order to make the sequence explicit to the user.
	fmt.Fprintln(cmd.OutOrStdout(), formatter.Success(uninstallMessage))
	fmt.Fprintln(cmd.OutOrStdout(), formatter.Success(purgeMessage))
	return nil
}

// removeInstallArtifacts deletes the backend binary and all known wrapper
// files from the user's local bin directory. It ignores missing files so the
// uninstall flow stays idempotent.
func removeInstallArtifacts() error {
	components := append([]string{install.BinaryName}, install.KnownWrapperNames()...)
	for _, component := range components {
		path, err := install.GetInstallPath(component)
		if err != nil {
			return err
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}
	return nil
}

// defaultConfigDir resolves the standard config directory without creating it.
// Purge needs the path so it can remove the directory tree in one step.
func defaultConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	return filepath.Join(home, ".config", config.ConfigDir), nil
}
