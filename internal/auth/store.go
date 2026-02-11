package auth

import "errors"

// Store defines the interface for per-device API key storage.
type Store interface {
	SetAPIKey(deviceName, key string) error
	GetAPIKey(deviceName string) (string, error)
	DeleteAPIKey(deviceName string) error
	HasAPIKey(deviceName string) (bool, error)
}

var (
	ErrNoAPIKey              = errors.New("no API key configured")
	errNoTTY                 = errors.New("no TTY available for keyring file backend password prompt")
	errInvalidKeyringBackend = errors.New("invalid keyring backend")
	errKeyringTimeout        = errors.New("keyring connection timed out")
	errEmptyAPIKey           = errors.New("API key cannot be empty")
	errEmptyDeviceName       = errors.New("device name cannot be empty")
)
