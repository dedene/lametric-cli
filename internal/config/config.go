package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Device holds connection details for a single LaMetric device.
type Device struct {
	IP string `yaml:"ip" json:"ip"`
}

// Output holds output preferences.
type Output struct {
	Color string `yaml:"color,omitempty" json:"color,omitempty"`
}

// File represents the lametric-cli YAML configuration.
type File struct {
	DefaultDevice  string            `yaml:"default_device,omitempty" json:"default_device,omitempty"`
	Devices        map[string]Device `yaml:"devices,omitempty" json:"devices,omitempty"`
	Output         Output            `yaml:"output,omitempty" json:"output,omitempty"`
	KeyringBackend string            `yaml:"keyring_backend,omitempty" json:"keyring_backend,omitempty"`
}

// DeviceByAlias returns the device for the given alias.
// If alias is empty, the default device is used.
func (f File) DeviceByAlias(alias string) (Device, error) {
	if alias == "" {
		alias = f.DefaultDevice
	}

	if alias == "" {
		return Device{}, fmt.Errorf("no device alias specified and no default_device configured")
	}

	d, ok := f.Devices[alias]
	if !ok {
		return Device{}, fmt.Errorf("device %q not found in config", alias)
	}

	return d, nil
}

// ConfigExists checks whether the config file exists on disk.
func ConfigExists() (bool, error) {
	path, err := ConfigPath()
	if err != nil {
		return false, err
	}

	return configExistsAt(path)
}

func configExistsAt(path string) (bool, error) {
	if _, statErr := os.Stat(path); statErr != nil {
		if os.IsNotExist(statErr) {
			return false, nil
		}

		return false, fmt.Errorf("stat config: %w", statErr)
	}

	return true, nil
}

// ReadConfig reads the YAML config file. Returns zero File{} if file does not exist.
func ReadConfig() (File, error) {
	path, err := ConfigPath()
	if err != nil {
		return File{}, err
	}

	return readConfigFrom(path)
}

func readConfigFrom(path string) (File, error) {
	b, err := os.ReadFile(path) //nolint:gosec // config file path
	if err != nil {
		if os.IsNotExist(err) {
			return File{}, nil
		}

		return File{}, fmt.Errorf("read config: %w", err)
	}

	var cfg File
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return File{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	return cfg, nil
}

// WriteConfig writes the YAML config file atomically using a .tmp + rename pattern.
func WriteConfig(cfg File) error {
	_, err := EnsureDir()
	if err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	return writeConfigTo(path, cfg)
}

func writeConfigTo(path string, cfg File) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config yaml: %w", err)
	}

	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("commit config: %w", err)
	}

	return nil
}
