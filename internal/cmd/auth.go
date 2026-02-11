package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/dedene/lametric-cli/internal/auth"
	"github.com/dedene/lametric-cli/internal/config"
)

// AuthCmd manages API key storage.
type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" name:"set-key" help:"Store API key for a device"`
	Status AuthStatusCmd `cmd:"" help:"Show API key status"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove stored API key"`
}

// AuthSetKeyCmd stores an API key in the system keyring.
type AuthSetKeyCmd struct {
	Stdin bool `help:"Read API key from stdin (for scripts)"`
}

// Run prompts for an API key and stores it for the resolved device.
func (c *AuthSetKeyCmd) Run(flags *RootFlags) error {
	deviceName, err := resolveDeviceName(flags.Device)
	if err != nil {
		return err
	}

	var key string

	if c.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			key = strings.TrimSpace(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read from stdin: %w", err)
		}
	} else {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("not a terminal; use --stdin flag to read from pipe")
		}

		fmt.Printf("Enter API key for %q: ", deviceName)

		bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			return fmt.Errorf("read API key: %w", err)
		}

		key = strings.TrimSpace(string(bytes))
	}

	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	store, err := auth.OpenDefault()
	if err != nil {
		return fmt.Errorf("open keyring: %w", err)
	}

	if err := store.SetAPIKey(deviceName, key); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	fmt.Fprintf(os.Stdout, "API key stored for %q.\n", deviceName)

	return nil
}

// AuthStatusCmd shows the current API key status.
type AuthStatusCmd struct{}

// Run checks env var, keyring, and reports status.
func (c *AuthStatusCmd) Run(flags *RootFlags) error {
	if os.Getenv("LAMETRIC_API_KEY") != "" {
		fmt.Fprintln(os.Stdout, "API key: set via LAMETRIC_API_KEY environment variable")

		return nil
	}

	backendInfo, err := auth.ResolveKeyringBackendInfo()
	if err != nil {
		return err
	}

	configPath, _ := config.ConfigPath()
	keyringDir, _ := config.KeyringDir()

	fmt.Fprintf(os.Stdout, "Config path:     %s\n", configPath)
	fmt.Fprintf(os.Stdout, "Keyring dir:     %s\n", keyringDir)
	fmt.Fprintf(os.Stdout, "Keyring backend: %s (source: %s)\n", backendInfo.Value, backendInfo.Source)

	deviceName, err := resolveDeviceName(flags.Device)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Device:          not configured (run: lametric setup)\n")

		return nil
	}

	fmt.Fprintf(os.Stdout, "Device:          %s\n", deviceName)

	store, err := auth.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stdout, "API key:         error opening keyring: %v\n", err)

		return nil
	}

	hasKey, err := store.HasAPIKey(deviceName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "API key:         error checking: %v\n", err)

		return nil
	}

	if hasKey {
		fmt.Fprintln(os.Stdout, "API key:         stored in system keyring")
	} else {
		fmt.Fprintln(os.Stdout, "API key:         not configured")
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "Run 'lametric setup' to configure your device.")
	}

	return nil
}

// AuthRemoveCmd removes the stored API key.
type AuthRemoveCmd struct{}

// Run deletes the API key from keyring for the resolved device.
func (c *AuthRemoveCmd) Run(flags *RootFlags) error {
	deviceName, err := resolveDeviceName(flags.Device)
	if err != nil {
		return err
	}

	store, err := auth.OpenDefault()
	if err != nil {
		return fmt.Errorf("open keyring: %w", err)
	}

	if err := store.DeleteAPIKey(deviceName); err != nil {
		if errors.Is(err, auth.ErrNoAPIKey) {
			fmt.Fprintf(os.Stdout, "No API key stored for %q.\n", deviceName)

			return nil
		}

		return fmt.Errorf("remove API key: %w", err)
	}

	fmt.Fprintf(os.Stdout, "API key removed for %q.\n", deviceName)

	return nil
}

// resolveDeviceName returns the device alias from the flag, env, or config default.
func resolveDeviceName(flagDevice string) (string, error) {
	if flagDevice != "" {
		return flagDevice, nil
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		return "", fmt.Errorf("read config: %w", err)
	}

	if cfg.DefaultDevice == "" {
		return "", fmt.Errorf("no device specified; use --device flag or run: lametric setup")
	}

	return cfg.DefaultDevice, nil
}
