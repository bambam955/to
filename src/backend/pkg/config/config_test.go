package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	installDir := filepath.Join(home, "nested", "bin")
	if err := Save(Config{InstallDir: installDir}); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.InstallDir != installDir {
		t.Fatalf("Load() = %+v, want InstallDir %q", loaded, installDir)
	}

	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() returned error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", path, err)
	}
}

func TestLoadMissingConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should return an error for a missing config file")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() error = %v, want os.ErrNotExist", err)
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("rejects empty install dir", func(t *testing.T) {
		err := Config{}.Validate()
		if err == nil {
			t.Fatal("Validate() should reject empty install_dir")
		}
		if !strings.Contains(err.Error(), "install_dir is empty") {
			t.Fatalf("Validate() error = %v, want empty install_dir error", err)
		}
	})

	t.Run("rejects relative install dir", func(t *testing.T) {
		err := Config{InstallDir: "relative/bin"}.Validate()
		if err == nil {
			t.Fatal("Validate() should reject relative install_dir")
		}
		if !strings.Contains(err.Error(), "absolute path") {
			t.Fatalf("Validate() error = %v, want absolute path error", err)
		}
	})
}

func TestLoadRejectsRelativeInstallDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".config", ConfigDir)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	path := filepath.Join(configDir, ConfigFileName)
	if err := os.WriteFile(path, []byte(`install_dir = "relative/bin"`+"\n"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should reject relative install_dir")
	}
	if !strings.Contains(err.Error(), "invalid config") {
		t.Fatalf("Load() error = %v, want invalid config error", err)
	}
}
