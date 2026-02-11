package cmd

import (
	"context"
	"fmt"

	"github.com/dedene/lametric-cli/internal/api"
)

// DeviceCmd shows device information.
type DeviceCmd struct{}

// Run fetches and displays device info.
func (c *DeviceCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var device api.Device
	if err := client.Get(context.Background(), "/api/v2/device", &device); err != nil {
		return fmt.Errorf("get device info: %w", err)
	}

	f := newFormatter(flags)

	return f.OutputSingle(device, [][2]string{
		{"Name", device.Name},
		{"Model", device.Model},
		{"Serial", device.SerialNumber},
		{"OS Version", device.OSVersion},
		{"Mode", device.Mode},
		{"WiFi IP", device.WiFi.IP},
		{"WiFi SSID", device.WiFi.SSID},
	})
}
