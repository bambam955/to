package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"to/pkg/config"
	"to/pkg/install"
)

// writeInstallConfig persists a canonical install path for uninstall and
// purge tests that need recorded metadata.
func writeInstallConfig(t *testing.T, home, installDir string) {
	t.Helper()

	t.Setenv("HOME", home)
	if err := config.Save(config.Config{InstallDir: installDir}); err != nil {
		t.Fatalf("failed to write install config: %v", err)
	}
}

// chdir swaps the working directory for a test and restores the original
// directory when the test ends.
func chdir(t *testing.T, dir string) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change working dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(original)
	})
}

func TestInstallManagementCommands(t *testing.T) {
	t.Run("uninstall removes backend and all wrappers", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		installDir := filepath.Join(home, install.TargetDir)
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			t.Fatalf("failed to create install dir: %v", err)
		}

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path, err := install.GetInstallPath(component)
			if err != nil {
				t.Fatalf("failed to resolve install path for %s: %v", component, err)
			}
			if err := os.WriteFile(path, []byte(component), 0o644); err != nil {
				t.Fatalf("failed to create install artifact %s: %v", path, err)
			}
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("y\n"))
		rootCmd.SetArgs([]string{"--uninstall"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path, err := install.GetInstallPath(component)
			if err != nil {
				t.Fatalf("failed to resolve install path for %s: %v", component, err)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("expected %s to be removed, stat err=%v", path, err)
			}
		}

		output := stdout.String()
		if !strings.Contains(output, "Confirm TO uninstall (y/n):") {
			t.Fatalf("expected uninstall confirmation prompt, got: %q", output)
		}
		if !strings.Contains(output, "removed TO backend and shell wrappers") {
			t.Fatalf("expected uninstall success message, got: %q", output)
		}
	})

	t.Run("uninstall cancels when declined", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("n\n"))
		rootCmd.SetArgs([]string{"-U"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "Confirm TO uninstall (y/n):") {
			t.Fatalf("expected uninstall confirmation prompt, got: %q", output)
		}
		if !strings.Contains(output, "cancelled") {
			t.Fatalf("expected cancellation message, got: %q", output)
		}
	})

	t.Run("uninstall respects install config", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		installDir := filepath.Join(home, "nested", "bin")
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			t.Fatalf("failed to create install dir: %v", err)
		}
		writeInstallConfig(t, home, installDir)
		t.Setenv("TO_INSTALL_DIR", filepath.Join(t.TempDir(), "wrong"))

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path := filepath.Join(installDir, component)
			if err := os.WriteFile(path, []byte(component), 0o644); err != nil {
				t.Fatalf("failed to create install artifact %s: %v", path, err)
			}
		}

		chdir(t, t.TempDir())

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("y\n"))
		rootCmd.SetArgs([]string{"-U"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path := filepath.Join(installDir, component)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("expected %s to be removed, stat err=%v", path, err)
			}
		}
	})

	t.Run("purge removes config data and install artifacts", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		installDir := filepath.Join(home, install.TargetDir)
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			t.Fatalf("failed to create install dir: %v", err)
		}
		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path, err := install.GetInstallPath(component)
			if err != nil {
				t.Fatalf("failed to resolve install path for %s: %v", component, err)
			}
			if err := os.WriteFile(path, []byte(component), 0o644); err != nil {
				t.Fatalf("failed to create install artifact %s: %v", path, err)
			}
		}

		configDir, err := config.GetConfigDir()
		if err != nil {
			t.Fatalf("failed to create config dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(configDir, "database.json"), []byte("{}"), 0o644); err != nil {
			t.Fatalf("failed to create config file: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("y\n"))
		rootCmd.SetArgs([]string{"--purge"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path, err := install.GetInstallPath(component)
			if err != nil {
				t.Fatalf("failed to resolve install path for %s: %v", component, err)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("expected %s to be removed, stat err=%v", path, err)
			}
		}

		if _, err := os.Stat(configDir); !os.IsNotExist(err) {
			t.Fatalf("expected config directory to be removed, stat err=%v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "Confirm TO purge (y/n):") {
			t.Fatalf("expected purge confirmation prompt, got: %q", output)
		}
		if !strings.Contains(output, "removed TO backend and shell wrappers") {
			t.Fatalf("expected uninstall message, got: %q", output)
		}
		if !strings.Contains(output, "purged TO database") {
			t.Fatalf("expected purge message, got: %q", output)
		}
		if strings.Index(output, "removed TO backend and shell wrappers") > strings.Index(output, "purged TO database") {
			t.Fatalf("expected uninstall message before purge message, got: %q", output)
		}
	})

	t.Run("purge respects install config", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		installDir := filepath.Join(home, "nested", "bin")
		if err := os.MkdirAll(installDir, 0o755); err != nil {
			t.Fatalf("failed to create install dir: %v", err)
		}
		writeInstallConfig(t, home, installDir)
		t.Setenv("TO_INSTALL_DIR", filepath.Join(t.TempDir(), "wrong"))

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path := filepath.Join(installDir, component)
			if err := os.WriteFile(path, []byte(component), 0o644); err != nil {
				t.Fatalf("failed to create install artifact %s: %v", path, err)
			}
		}

		chdir(t, t.TempDir())

		configDir, err := config.GetConfigDir()
		if err != nil {
			t.Fatalf("failed to create config dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(configDir, "database.json"), []byte("{}"), 0o644); err != nil {
			t.Fatalf("failed to create config file: %v", err)
		}

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("y\n"))
		rootCmd.SetArgs([]string{"-P"})

		err = rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		for _, component := range append([]string{install.BinaryName}, install.KnownWrapperNames()...) {
			path := filepath.Join(installDir, component)
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("expected %s to be removed, stat err=%v", path, err)
			}
		}

		if _, err := os.Stat(configDir); !os.IsNotExist(err) {
			t.Fatalf("expected config directory to be removed, stat err=%v", err)
		}
	})

	t.Run("purge cancels when declined", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetIn(strings.NewReader("n\n"))
		rootCmd.SetArgs([]string{"-P"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "Confirm TO purge (y/n):") {
			t.Fatalf("expected purge confirmation prompt, got: %q", output)
		}
		if !strings.Contains(output, "cancelled") {
			t.Fatalf("expected cancellation message, got: %q", output)
		}
	})

	t.Run("conflicting operation flags fail fast", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		rootCmd.SetArgs([]string{"--uninstall", "--purge"})

		err := rootCmd.Execute()
		if err == nil {
			t.Fatal("expected error for conflicting flags, got nil")
		}
		if !strings.Contains(err.Error(), "only one operation flag") {
			t.Fatalf("expected conflict error, got: %v", err)
		}
	})
}
