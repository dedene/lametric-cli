package cmd

import "fmt"

// DiscoverCmd finds LaMetric devices on the local network.
type DiscoverCmd struct{}

// Run prints a placeholder message until discovery is implemented.
func (c *DiscoverCmd) Run() error {
	fmt.Println("Device discovery not yet implemented")
	fmt.Println("Use: lametric setup --ip=<device-ip>")

	return nil
}
