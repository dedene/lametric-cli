package auth

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"golang.org/x/term"

	"github.com/dedene/lametric-cli/internal/config"
)

// KeyringStore implements Store using the system keyring.
type KeyringStore struct {
	ring keyring.Keyring
}

const (
	keyringPrefix      = "lametric:"
	keyringPasswordEnv = "LAMETRIC_KEYRING_PASSWORD" //nolint:gosec // env var name
	keyringBackendEnv  = "LAMETRIC_KEYRING_BACKEND"  //nolint:gosec // env var name
	apiKeyEnv          = "LAMETRIC_API_KEY"           //nolint:gosec // env var name
)

var (
	openKeyringFunc = openKeyring
	keyringOpenFunc = keyring.Open
)

// KeyringBackendInfo holds the resolved keyring backend value and its source.
type KeyringBackendInfo struct {
	Value  string
	Source string
}

const (
	keyringBackendSourceEnv     = "env"
	keyringBackendSourceConfig  = "config"
	keyringBackendSourceDefault = "default"
	keyringBackendAuto          = "auto"
	keyringOpenTimeout          = 5 * time.Second
)

// ResolveKeyringBackendInfo determines the keyring backend from env, config, or default.
func ResolveKeyringBackendInfo() (KeyringBackendInfo, error) {
	if v := normalizeKeyringBackend(os.Getenv(keyringBackendEnv)); v != "" {
		return KeyringBackendInfo{Value: v, Source: keyringBackendSourceEnv}, nil
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return KeyringBackendInfo{}, fmt.Errorf("resolve keyring backend: %w", err)
	}

	if cfg.KeyringBackend != "" {
		if v := normalizeKeyringBackend(cfg.KeyringBackend); v != "" {
			return KeyringBackendInfo{Value: v, Source: keyringBackendSourceConfig}, nil
		}
	}

	return KeyringBackendInfo{Value: keyringBackendAuto, Source: keyringBackendSourceDefault}, nil
}

func allowedBackends(info KeyringBackendInfo) ([]keyring.BackendType, error) {
	switch info.Value {
	case "", keyringBackendAuto:
		return nil, nil
	case "keychain":
		return []keyring.BackendType{keyring.KeychainBackend}, nil
	case "file":
		return []keyring.BackendType{keyring.FileBackend}, nil
	default:
		return nil, fmt.Errorf("%w: %q (expected %s, keychain, or file)", errInvalidKeyringBackend, info.Value, keyringBackendAuto)
	}
}

func wrapKeychainError(err error) error {
	if err == nil {
		return nil
	}

	if isKeychainLockedError(err.Error()) {
		return fmt.Errorf("%w\n\nYour macOS keychain is locked. To unlock it, run:\n  security unlock-keychain ~/Library/Keychains/login.keychain-db", err)
	}

	return err
}

func isKeychainLockedError(msg string) bool {
	return strings.Contains(msg, "errSecInteractionNotAllowed") ||
		strings.Contains(msg, "The user name or passphrase you entered is not correct")
}

func fileKeyringPasswordFunc() keyring.PromptFunc {
	password := os.Getenv(keyringPasswordEnv)
	if password != "" {
		return keyring.FixedStringPrompt(password)
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		return keyring.TerminalPrompt
	}

	return func(_ string) (string, error) {
		return "", fmt.Errorf("%w; set %s", errNoTTY, keyringPasswordEnv)
	}
}

func normalizeKeyringBackend(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func shouldForceFileBackend(goos string, backendInfo KeyringBackendInfo, dbusAddr string) bool {
	return goos == "linux" && backendInfo.Value == keyringBackendAuto && dbusAddr == ""
}

func shouldUseKeyringTimeout(goos string, backendInfo KeyringBackendInfo, dbusAddr string) bool {
	return goos == "linux" && backendInfo.Value == keyringBackendAuto && dbusAddr != ""
}

func deviceKeyringKey(deviceName string) string {
	return keyringPrefix + deviceName
}

func openKeyring() (keyring.Keyring, error) {
	keyringDir, err := config.EnsureKeyringDir()
	if err != nil {
		return nil, fmt.Errorf("ensure keyring dir: %w", err)
	}

	backendInfo, err := ResolveKeyringBackendInfo()
	if err != nil {
		return nil, err
	}

	backends, err := allowedBackends(backendInfo)
	if err != nil {
		return nil, err
	}

	dbusAddr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")

	if shouldForceFileBackend(runtime.GOOS, backendInfo, dbusAddr) {
		backends = []keyring.BackendType{keyring.FileBackend}
	}

	cfg := keyring.Config{
		ServiceName:              config.AppName,
		KeychainTrustApplication: false,
		AllowedBackends:          backends,
		FileDir:                  keyringDir,
		FilePasswordFunc:         fileKeyringPasswordFunc(),
	}

	if shouldUseKeyringTimeout(runtime.GOOS, backendInfo, dbusAddr) {
		return openKeyringWithTimeout(cfg, keyringOpenTimeout)
	}

	ring, err := keyringOpenFunc(cfg)
	if err != nil {
		return nil, fmt.Errorf("open keyring: %w", err)
	}

	return ring, nil
}

type keyringResult struct {
	ring keyring.Keyring
	err  error
}

func openKeyringWithTimeout(cfg keyring.Config, timeout time.Duration) (keyring.Keyring, error) {
	ch := make(chan keyringResult, 1)

	go func() {
		ring, err := keyringOpenFunc(cfg)
		ch <- keyringResult{ring, err}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			return nil, fmt.Errorf("open keyring: %w", res.err)
		}

		return res.ring, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("%w after %v (D-Bus SecretService may be unresponsive); "+
			"set LAMETRIC_KEYRING_BACKEND=file and LAMETRIC_KEYRING_PASSWORD=<password> to use encrypted file storage instead",
			errKeyringTimeout, timeout)
	}
}

// OpenDefault opens the default keyring store.
func OpenDefault() (Store, error) {
	ring, err := openKeyringFunc()
	if err != nil {
		return nil, err
	}

	return &KeyringStore{ring: ring}, nil
}

// SetAPIKey stores an API key for a device in the keyring.
func (s *KeyringStore) SetAPIKey(deviceName, key string) error {
	deviceName = strings.TrimSpace(deviceName)
	if deviceName == "" {
		return errEmptyDeviceName
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return errEmptyAPIKey
	}

	if err := s.ring.Set(keyring.Item{
		Key:  deviceKeyringKey(deviceName),
		Data: []byte(key),
	}); err != nil {
		return wrapKeychainError(fmt.Errorf("store API key for %q: %w", deviceName, err))
	}

	return nil
}

// GetAPIKey retrieves the API key for a device from the keyring.
func (s *KeyringStore) GetAPIKey(deviceName string) (string, error) {
	deviceName = strings.TrimSpace(deviceName)
	if deviceName == "" {
		return "", errEmptyDeviceName
	}

	item, err := s.ring.Get(deviceKeyringKey(deviceName))
	if err != nil {
		if isKeyNotFound(err) {
			return "", ErrNoAPIKey
		}

		return "", fmt.Errorf("read API key for %q: %w", deviceName, err)
	}

	return string(item.Data), nil
}

// DeleteAPIKey removes the API key for a device from the keyring.
func (s *KeyringStore) DeleteAPIKey(deviceName string) error {
	deviceName = strings.TrimSpace(deviceName)
	if deviceName == "" {
		return errEmptyDeviceName
	}

	if err := s.ring.Remove(deviceKeyringKey(deviceName)); err != nil {
		if isKeyNotFound(err) {
			return ErrNoAPIKey
		}

		return fmt.Errorf("delete API key for %q: %w", deviceName, err)
	}

	return nil
}

// HasAPIKey checks whether an API key is stored for a device.
func (s *KeyringStore) HasAPIKey(deviceName string) (bool, error) {
	deviceName = strings.TrimSpace(deviceName)
	if deviceName == "" {
		return false, errEmptyDeviceName
	}

	_, err := s.ring.Get(deviceKeyringKey(deviceName))
	if err != nil {
		if isKeyNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("check API key for %q: %w", deviceName, err)
	}

	return true, nil
}

func isKeyNotFound(err error) bool {
	return err == keyring.ErrKeyNotFound
}

// GetAPIKey checks the env var first, then falls back to the keyring for the given device.
func GetAPIKey(deviceName string) (string, error) {
	if envKey := os.Getenv(apiKeyEnv); envKey != "" {
		return envKey, nil
	}

	store, err := OpenDefault()
	if err != nil {
		return "", err
	}

	key, err := store.GetAPIKey(deviceName)
	if err != nil {
		return "", fmt.Errorf("get API key: %w", err)
	}

	return key, nil
}
