package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dedene/lametric-cli/internal/api"
)

// WiFiCmd shows wifi status (read-only).
type WiFiCmd struct{}

func (c *WiFiCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var wifi api.WiFi
	if err := client.Get(context.Background(), "/api/v2/device/wifi", &wifi); err != nil {
		return fmt.Errorf("get wifi: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(wifi, [][2]string{
		{"Active", strconv.FormatBool(wifi.Active)},
		{"SSID", wifi.SSID},
		{"IP", wifi.IP},
		{"MAC", wifi.MAC},
		{"Netmask", wifi.Netmask},
		{"Strength", strconv.Itoa(wifi.Strength)},
		{"Encryption", wifi.Encryption},
	})
}
