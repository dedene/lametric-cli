package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/dedene/lametric-cli/internal/discovery"
)

// DiscoverCmd finds LaMetric devices on the local network.
type DiscoverCmd struct {
	Timeout time.Duration `help:"Discovery timeout" default:"10s" short:"t"`
}

// Run discovers devices and displays them.
func (c *DiscoverCmd) Run(flags *RootFlags) error {
	f := newFormatter(flags)

	devices, err := discovery.Discover(context.Background(), c.Timeout)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}

	if len(devices) == 0 {
		fmt.Println("No LaMetric devices found on the network.")
		return nil
	}

	rows := make([][]string, len(devices))
	for i, d := range devices {
		rows[i] = []string{d.Name, d.IP, fmt.Sprintf("%d", d.Port), d.Model}
	}

	return f.Output(devices, []string{"NAME", "IP", "PORT", "MODEL"}, rows)
}
