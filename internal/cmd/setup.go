package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/dedene/lametric-cli/internal/auth"
	"github.com/dedene/lametric-cli/internal/config"
)

// SetupCmd is the interactive setup wizard for configuring a LaMetric device.
type SetupCmd struct {
	IP   string `help:"Device IP address" short:"i"`
	Name string `help:"Device alias" short:"n"`
}

// Run walks the user through device setup: IP, API key, config.
func (c *SetupCmd) Run() error {
	ip := c.IP
	name := c.Name

	reader := bufio.NewReader(os.Stdin)

	// Prompt for IP if not provided.
	if ip == "" {
		fmt.Print("Device IP address: ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read IP: %w", err)
		}

		ip = strings.TrimSpace(line)
	}

	if ip == "" {
		return fmt.Errorf("device IP is required")
	}

	// Prompt for device alias if not provided.
	if name == "" {
		fmt.Printf("Device alias [%s]: ", ip)

		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read alias: %w", err)
		}

		name = strings.TrimSpace(line)
		if name == "" {
			name = ip
		}
	}

	// Prompt for API key (masked).
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("not a terminal; setup requires interactive input")
	}

	fmt.Print("API key (input hidden): ")

	keyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		return fmt.Errorf("read API key: %w", err)
	}

	key := strings.TrimSpace(string(keyBytes))
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Store API key in keyring.
	store, err := auth.OpenDefault()
	if err != nil {
		return fmt.Errorf("open keyring: %w", err)
	}

	if err := store.SetAPIKey(name, key); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	// Save device in config.
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if cfg.Devices == nil {
		cfg.Devices = make(map[string]config.Device)
	}

	cfg.Devices[name] = config.Device{IP: ip}

	// Set as default if first device or no default.
	if cfg.DefaultDevice == "" {
		cfg.DefaultDevice = name
	}

	if err := config.WriteConfig(cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println()
	fmt.Printf("Device %q configured successfully.\n", name)
	fmt.Printf("  IP:      %s\n", ip)
	fmt.Printf("  Default: %v\n", cfg.DefaultDevice == name)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  lametric device          Show device info")
	fmt.Println("  lametric notify \"Hello\"  Send a notification")

	return nil
}
