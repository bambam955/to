package install

import (
	"os"
	"path/filepath"
	"testing"

	"to/pkg/config"
)

// writeRawInstallConfig writes a test fixture directly so we can exercise
// fallback behavior against malformed config.toml contents.
func writeRawInstallConfig(t *testing.T, home, contents string) {
	t.Helper()

	configDir := filepath.Join(home, ".config", config.ConfigDir)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, config.ConfigFileName)
	if err := os.WriteFile(configPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
}

func TestGetInstallPath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name      string
		component string
		expected  string
	}{
		{
			name:      "binary component",
			component: BinaryName,
			expected:  filepath.Join(home, TargetDir, BinaryName),
		},
		{
			name:      "wrapper component",
			component: WrapperName,
			expected:  filepath.Join(home, TargetDir, WrapperName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetInstallPath(tt.component)
			if err != nil {
				t.Errorf("GetInstallPath(%q) returned error: %v", tt.component, err)
			}
			if result != tt.expected {
				t.Errorf("GetInstallPath(%q) = %q, want %q", tt.component, result, tt.expected)
			}
		})
	}

	t.Run("install dir override", func(t *testing.T) {
		overrideDir := filepath.Join(t.TempDir(), "custom install")
		t.Setenv("TO_INSTALL_DIR", overrideDir)

		result, err := GetInstallPath(BinaryName)
		if err != nil {
			t.Fatalf("GetInstallPath(%q) returned error: %v", BinaryName, err)
		}

		expected := filepath.Join(overrideDir, BinaryName)
		if result != expected {
			t.Errorf("GetInstallPath(%q) = %q, want %q", BinaryName, result, expected)
		}
	})
}

func TestEnsureInstallDir(t *testing.T) {
	// Create a temporary home directory for testing
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	err := EnsureInstallDir()
	if err != nil {
		t.Fatalf("EnsureInstallDir() returned error: %v", err)
	}

	expectedDir := filepath.Join(tmpHome, TargetDir)
	info, err := os.Stat(expectedDir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("Path exists but is not a directory")
	}

	// Verify it's idempotent - should not fail if already exists
	err = EnsureInstallDir()
	if err != nil {
		t.Errorf("Second call to EnsureInstallDir() returned error: %v", err)
	}

	t.Run("respects TO_INSTALL_DIR override", func(t *testing.T) {
		overrideDir := filepath.Join(t.TempDir(), "custom install")
		t.Setenv("TO_INSTALL_DIR", overrideDir)

		if err := EnsureInstallDir(); err != nil {
			t.Fatalf("EnsureInstallDir() returned error: %v", err)
		}

		info, err := os.Stat(overrideDir)
		if err != nil {
			t.Fatalf("Override directory was not created: %v", err)
		}
		if !info.IsDir() {
			t.Errorf("Path exists but is not a directory")
		}
	})
}

func TestVerifyPath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	installDir := filepath.Join(home, TargetDir)

	// Save original PATH
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Test when install dir is in PATH
	os.Setenv("PATH", origPath+string(os.PathListSeparator)+installDir)
	inPath, err := VerifyPath()
	if err != nil {
		t.Fatalf("VerifyPath() returned error: %v", err)
	}
	if !inPath {
		t.Errorf("VerifyPath() = %v, want %v", inPath, true)
	}

	// Test when install dir is not in PATH
	os.Setenv("PATH", "/usr/bin:/bin")
	inPath, err = VerifyPath()
	if err != nil {
		t.Fatalf("VerifyPath() returned error: %v", err)
	}
	if inPath {
		t.Errorf("VerifyPath() = %v, want %v", inPath, false)
	}

	t.Run("with TO_INSTALL_DIR override", func(t *testing.T) {
		overrideDir := filepath.Join(t.TempDir(), "custom install")
		t.Setenv("TO_INSTALL_DIR", overrideDir)

		origPath := os.Getenv("PATH")
		t.Setenv("PATH", origPath+string(os.PathListSeparator)+overrideDir)
		inPath, err := VerifyPath()
		if err != nil {
			t.Fatalf("VerifyPath() returned error: %v", err)
		}
		if !inPath {
			t.Errorf("VerifyPath() = %v, want %v", inPath, true)
		}
	})

	t.Run("prefers recorded install dir", func(t *testing.T) {
		recordedDir := filepath.Join(t.TempDir(), "recorded install")
		if err := os.MkdirAll(recordedDir, 0o755); err != nil {
			t.Fatalf("failed to create recorded install dir: %v", err)
		}
		if err := config.Save(config.Config{InstallDir: recordedDir}); err != nil {
			t.Fatalf("failed to write install config: %v", err)
		}

		t.Setenv("TO_INSTALL_DIR", filepath.Join(t.TempDir(), "different"))
		t.Setenv("PATH", "/usr/bin:/bin")

		inPath, err := VerifyPath()
		if err != nil {
			t.Fatalf("VerifyPath() returned error: %v", err)
		}
		if inPath {
			t.Errorf("VerifyPath() = %v, want %v", inPath, false)
		}

		t.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+recordedDir)
		inPath, err = VerifyPath()
		if err != nil {
			t.Fatalf("VerifyPath() returned error: %v", err)
		}
		if !inPath {
			t.Errorf("VerifyPath() = %v, want %v", inPath, true)
		}
	})
}

func TestResolveInstallDirFallsBackWhenConfigIsBroken(t *testing.T) {
	t.Run("uses TO_INSTALL_DIR when config is invalid", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		writeRawInstallConfig(t, home, `install_dir = "relative/bin"`+"\n")
		overrideDir := filepath.Join(t.TempDir(), "custom install")
		t.Setenv("TO_INSTALL_DIR", overrideDir)

		installDir, err := ResolveInstallDir()
		if err != nil {
			t.Fatalf("ResolveInstallDir() returned error: %v", err)
		}
		if installDir != overrideDir {
			t.Fatalf("ResolveInstallDir() = %q, want %q", installDir, overrideDir)
		}
	})

	t.Run("uses default install dir when config is malformed", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		writeRawInstallConfig(t, home, `install_dir = [`+"\n")
		t.Setenv("TO_INSTALL_DIR", "")

		installDir, err := ResolveInstallDir()
		if err != nil {
			t.Fatalf("ResolveInstallDir() returned error: %v", err)
		}

		expected := filepath.Join(home, TargetDir)
		if installDir != expected {
			t.Fatalf("ResolveInstallDir() = %q, want %q", installDir, expected)
		}
	})
}
