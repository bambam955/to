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
		if !strings.Contains(output, "removed installed backend and shell wrappers") {
			t.Fatalf("expected uninstall success message, got: %q", output)
		}
	})

	t.Run("uninstall is a no-op when artifacts are missing", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"-U"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "removed installed backend and shell wrappers") {
			t.Fatalf("expected uninstall success message, got: %q", output)
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
		if !strings.Contains(output, "purged install artifacts and default config data") {
			t.Fatalf("expected purge success message, got: %q", output)
		}
	})

	t.Run("purge is a no-op when config data is missing", func(t *testing.T) {
		resetFlags(t)

		home := t.TempDir()
		t.Setenv("HOME", home)

		var stdout bytes.Buffer
		rootCmd.SetOut(&stdout)
		rootCmd.SetArgs([]string{"-P"})

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		output := stdout.String()
		if !strings.Contains(output, "purged install artifacts and default config data") {
			t.Fatalf("expected purge success message, got: %q", output)
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
