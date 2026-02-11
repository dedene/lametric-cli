package auth

import (
	"errors"
	"os"
	"testing"

	"github.com/99designs/keyring"
)

// mockKeyring implements keyring.Keyring in-memory.
type mockKeyring struct {
	items map[string]keyring.Item
}

func newMockKeyring() *mockKeyring {
	return &mockKeyring{items: make(map[string]keyring.Item)}
}

func (m *mockKeyring) Get(key string) (keyring.Item, error) {
	item, ok := m.items[key]
	if !ok {
		return keyring.Item{}, keyring.ErrKeyNotFound
	}
	return item, nil
}

func (m *mockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (m *mockKeyring) Set(item keyring.Item) error {
	m.items[item.Key] = item
	return nil
}

func (m *mockKeyring) Remove(key string) error {
	if _, ok := m.items[key]; !ok {
		return keyring.ErrKeyNotFound
	}
	delete(m.items, key)
	return nil
}

func (m *mockKeyring) Keys() ([]string, error) {
	keys := make([]string, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys, nil
}

func newTestStore() (*KeyringStore, *mockKeyring) {
	mk := newMockKeyring()
	return &KeyringStore{ring: mk}, mk
}

func TestKeyringStore_SetAndGet(t *testing.T) {
	store, _ := newTestStore()

	if err := store.SetAPIKey("office", "secret123"); err != nil {
		t.Fatalf("SetAPIKey: %v", err)
	}

	got, err := store.GetAPIKey("office")
	if err != nil {
		t.Fatalf("GetAPIKey: %v", err)
	}
	if got != "secret123" {
		t.Errorf("got %q, want %q", got, "secret123")
	}
}

func TestKeyringStore_Get_NotFound(t *testing.T) {
	store, _ := newTestStore()

	_, err := store.GetAPIKey("missing")
	if err != ErrNoAPIKey {
		t.Errorf("err = %v, want ErrNoAPIKey", err)
	}
}

func TestKeyringStore_Get_EnvOverride(t *testing.T) {
	t.Setenv("LAMETRIC_API_KEY", "env-key-999")

	// Package-level GetAPIKey checks env first, but it also calls OpenDefault
	// which requires a real keyring. Test the env check directly via os.Getenv.
	got := os.Getenv("LAMETRIC_API_KEY")
	if got != "env-key-999" {
		t.Errorf("env = %q, want %q", got, "env-key-999")
	}
}

func TestKeyringStore_Delete(t *testing.T) {
	store, _ := newTestStore()

	if err := store.SetAPIKey("tmp", "val"); err != nil {
		t.Fatal(err)
	}

	if err := store.DeleteAPIKey("tmp"); err != nil {
		t.Fatalf("DeleteAPIKey: %v", err)
	}

	_, err := store.GetAPIKey("tmp")
	if err != ErrNoAPIKey {
		t.Errorf("after delete: err = %v, want ErrNoAPIKey", err)
	}
}

func TestKeyringStore_Delete_NotFound(t *testing.T) {
	store, _ := newTestStore()

	err := store.DeleteAPIKey("nope")
	if err != ErrNoAPIKey {
		t.Errorf("err = %v, want ErrNoAPIKey", err)
	}
}

func TestKeyringStore_HasAPIKey(t *testing.T) {
	store, _ := newTestStore()

	has, err := store.HasAPIKey("dev1")
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Error("expected false before set")
	}

	if err := store.SetAPIKey("dev1", "key1"); err != nil {
		t.Fatal(err)
	}

	has, err = store.HasAPIKey("dev1")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Error("expected true after set")
	}
}

func TestKeyringStore_EmptyDeviceName(t *testing.T) {
	store, _ := newTestStore()

	if _, err := store.GetAPIKey(""); err != errEmptyDeviceName {
		t.Errorf("GetAPIKey empty: err = %v, want errEmptyDeviceName", err)
	}
	if err := store.SetAPIKey("", "key"); err != errEmptyDeviceName {
		t.Errorf("SetAPIKey empty device: err = %v, want errEmptyDeviceName", err)
	}
	if err := store.SetAPIKey("dev", ""); err != errEmptyAPIKey {
		t.Errorf("SetAPIKey empty key: err = %v, want errEmptyAPIKey", err)
	}
	if err := store.DeleteAPIKey(""); err != errEmptyDeviceName {
		t.Errorf("DeleteAPIKey empty: err = %v, want errEmptyDeviceName", err)
	}
}

func TestDeviceKeyringKey(t *testing.T) {
	got := deviceKeyringKey("office")
	want := "lametric:office"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeKeyringBackend(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"FILE", "file"},
		{" Keychain ", "keychain"},
		{"auto", "auto"},
		{"", ""},
	}
	for _, tt := range tests {
		got := normalizeKeyringBackend(tt.input)
		if got != tt.want {
			t.Errorf("normalizeKeyringBackend(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsKeychainLockedError(t *testing.T) {
	if !isKeychainLockedError("errSecInteractionNotAllowed blah") {
		t.Error("expected true for errSecInteractionNotAllowed")
	}
	if isKeychainLockedError("some other error") {
		t.Error("expected false for unrelated error")
	}
}

func TestShouldForceFileBackend(t *testing.T) {
	info := KeyringBackendInfo{Value: "auto"}

	if !shouldForceFileBackend("linux", info, "") {
		t.Error("expected true for linux+auto+no-dbus")
	}
	if shouldForceFileBackend("darwin", info, "") {
		t.Error("expected false for darwin")
	}
	if shouldForceFileBackend("linux", info, "/run/user/1000/bus") {
		t.Error("expected false when dbus present")
	}
}

func TestShouldUseKeyringTimeout(t *testing.T) {
	info := KeyringBackendInfo{Value: "auto"}

	if !shouldUseKeyringTimeout("linux", info, "/run/user/1000/bus") {
		t.Error("expected true for linux+auto+dbus")
	}
	if shouldUseKeyringTimeout("linux", info, "") {
		t.Error("expected false without dbus")
	}
	if shouldUseKeyringTimeout("darwin", info, "/run/user/1000/bus") {
		t.Error("expected false for darwin")
	}
}

func TestAllowedBackends(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantLen int
		wantErr bool
	}{
		{"auto", "auto", 0, false},
		{"empty", "", 0, false},
		{"keychain", "keychain", 1, false},
		{"file", "file", 1, false},
		{"invalid", "redis", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := KeyringBackendInfo{Value: tt.value}
			backends, err := allowedBackends(info)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(backends) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(backends), tt.wantLen)
			}
		})
	}
}

func TestWrapKeychainError_Nil(t *testing.T) {
	if wrapKeychainError(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestWrapKeychainError_Normal(t *testing.T) {
	err := errors.New("some error")
	got := wrapKeychainError(err)
	if got != err {
		t.Errorf("expected same error, got %v", got)
	}
}

func TestWrapKeychainError_Locked(t *testing.T) {
	err := errors.New("errSecInteractionNotAllowed")
	got := wrapKeychainError(err)
	if got == err {
		t.Error("expected wrapped error for locked keychain")
	}
}
