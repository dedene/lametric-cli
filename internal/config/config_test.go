package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfigFrom_FileNotExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg, err := readConfigFrom(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.DefaultDevice != "" {
		t.Errorf("expected empty DefaultDevice, got %q", cfg.DefaultDevice)
	}
}

func TestReadConfigFrom_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	data := []byte(`default_device: office
devices:
  office:
    ip: "192.168.1.42"
  bedroom:
    ip: "192.168.1.43"
output:
  color: auto
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	cfg, err := readConfigFrom(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.DefaultDevice != "office" {
		t.Errorf("DefaultDevice = %q, want %q", cfg.DefaultDevice, "office")
	}
	if len(cfg.Devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(cfg.Devices))
	}
	if cfg.Devices["office"].IP != "192.168.1.42" {
		t.Errorf("office IP = %q, want %q", cfg.Devices["office"].IP, "192.168.1.42")
	}
}

func TestWriteConfigTo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := File{
		DefaultDevice: "desk",
		Devices: map[string]Device{
			"desk": {IP: "10.0.0.5"},
		},
	}

	if err := writeConfigTo(path, cfg); err != nil {
		t.Fatalf("writeConfigTo: %v", err)
	}

	got, err := readConfigFrom(path)
	if err != nil {
		t.Fatalf("readConfigFrom: %v", err)
	}
	if got.DefaultDevice != "desk" {
		t.Errorf("DefaultDevice = %q, want %q", got.DefaultDevice, "desk")
	}
	if got.Devices["desk"].IP != "10.0.0.5" {
		t.Errorf("desk IP = %q, want %q", got.Devices["desk"].IP, "10.0.0.5")
	}
}

func TestDeviceByAlias(t *testing.T) {
	cfg := File{
		DefaultDevice: "main",
		Devices: map[string]Device{
			"main":   {IP: "10.0.0.1"},
			"backup": {IP: "10.0.0.2"},
		},
	}

	tests := []struct {
		name    string
		alias   string
		wantIP  string
		wantErr bool
	}{
		{"explicit alias", "backup", "10.0.0.2", false},
		{"empty uses default", "", "10.0.0.1", false},
		{"not found", "missing", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := cfg.DeviceByAlias(tt.alias)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && d.IP != tt.wantIP {
				t.Errorf("IP = %q, want %q", d.IP, tt.wantIP)
			}
		})
	}
}

func TestDeviceByAlias_NoDefault(t *testing.T) {
	cfg := File{
		Devices: map[string]Device{
			"main": {IP: "10.0.0.1"},
		},
	}

	_, err := cfg.DeviceByAlias("")
	if err == nil {
		t.Fatal("expected error for empty alias with no default")
	}
}

func TestConfigExistsAt(t *testing.T) {
	dir := t.TempDir()

	// File does not exist.
	exists, err := configExistsAt(filepath.Join(dir, "nope.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected false for missing file")
	}

	// File exists.
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	exists, err = configExistsAt(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected true for existing file")
	}
}

func TestReadConfigFrom_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(":::not yaml"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := readConfigFrom(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestWriteConfigTo_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := File{
		DefaultDevice: "test",
		Devices: map[string]Device{
			"test": {IP: "1.2.3.4"},
		},
		Output:         Output{Color: "auto"},
		KeyringBackend: "file",
	}

	if err := writeConfigTo(path, cfg); err != nil {
		t.Fatal(err)
	}

	got, err := readConfigFrom(path)
	if err != nil {
		t.Fatal(err)
	}

	if got.DefaultDevice != cfg.DefaultDevice {
		t.Errorf("DefaultDevice = %q, want %q", got.DefaultDevice, cfg.DefaultDevice)
	}
	if got.Output.Color != "auto" {
		t.Errorf("Output.Color = %q, want %q", got.Output.Color, "auto")
	}
	if got.KeyringBackend != "file" {
		t.Errorf("KeyringBackend = %q, want %q", got.KeyringBackend, "file")
	}
}

func TestDir(t *testing.T) {
	dir, err := Dir()
	if err != nil {
		t.Fatalf("Dir() error: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
}

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error: %v", err)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}
}

func TestKeyringDir(t *testing.T) {
	dir, err := KeyringDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Error("expected non-empty")
	}
}
