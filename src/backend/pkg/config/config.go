// Package config provides configuration and XDG Base Directory support.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"to/pkg/errors"
)

// Config stores app-wide settings persisted on disk.
type Config struct {
	InstallDir string `toml:"install_dir"`
}

// Validate checks whether the config contains values we can safely use.
// InstallDir must be a non-empty absolute path because uninstall and verify
// flows rely on it being canonical and independent from the current cwd.
func (c Config) Validate() error {
	if strings.TrimSpace(c.InstallDir) == "" {
		return errors.InvalidInput("install_dir is empty")
	}
	if !filepath.IsAbs(c.InstallDir) {
		return errors.InvalidInput(fmt.Sprintf("install_dir must be an absolute path: %s", c.InstallDir))
	}

	return nil
}

// ConfigPath returns the path to the persisted config.toml file.
// It does not create directories so callers can distinguish a missing file
// from a missing config directory.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(errors.ErrorTypeInternal, "cannot determine home directory", err)
	}

	return filepath.Join(home, ".config", ConfigDir, ConfigFileName), nil
}

// Load reads the persisted config from disk.
func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, err
		}
		return Config{}, errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot read config %s", path), err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, errors.Wrap(errors.ErrorTypeCorrupted, fmt.Sprintf("cannot parse config %s", path), err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, errors.Wrap(errors.ErrorTypeInvalid, fmt.Sprintf("invalid config %s", path), err)
	}

	return cfg, nil
}

// Save writes the config to disk using the same TOML shape the installer uses.
func Save(cfg Config) (err error) {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if _, err := GetConfigDir(); err != nil {
		return err
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot write config %s", path), err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = errors.Wrap(errors.ErrorTypePermission, fmt.Sprintf("cannot close config %s", path), closeErr)
		}
	}()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return errors.Wrap(errors.ErrorTypeCorrupted, fmt.Sprintf("cannot encode config %s", path), err)
	}

	return nil
}
