package cmd

import (
	"fmt"
	"os"

	"github.com/dedene/lametric-cli/internal/api"
	"github.com/dedene/lametric-cli/internal/auth"
	"github.com/dedene/lametric-cli/internal/config"
	"github.com/dedene/lametric-cli/internal/output"
)

// resolveClient resolves a device from config/flags and returns an API client.
func resolveClient(flags *RootFlags) (*api.Client, string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return nil, "", fmt.Errorf("read config: %w", err)
	}

	alias := flags.Device
	if alias == "" {
		alias = cfg.DefaultDevice
	}

	if alias == "" {
		return nil, "", fmt.Errorf("no device specified; use --device or set default_device in config")
	}

	dev, err := cfg.DeviceByAlias(alias)
	if err != nil {
		return nil, "", err
	}

	apiKey, err := auth.GetAPIKey(alias)
	if err != nil {
		return nil, "", fmt.Errorf("get API key for %q: %w", alias, err)
	}

	return api.NewClient(dev.IP, apiKey), alias, nil
}

// newFormatter creates a Formatter from root flags.
func newFormatter(flags *RootFlags) *output.Formatter {
	return output.NewFormatter(os.Stdout, flags.JSON, flags.Plain, flags.NoColor)
}
