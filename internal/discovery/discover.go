package discovery

import (
	"context"
	"fmt"
	"time"
)

// Device represents a discovered LaMetric device.
type Device struct {
	Name  string `json:"name"`
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Model string `json:"model"`
}

// String returns a human-readable representation of the device.
func (d Device) String() string {
	return fmt.Sprintf("%s (%s:%d)", d.Name, d.IP, d.Port)
}

// Discover finds LaMetric devices on the local network.
// Runs mDNS and SSDP in parallel, deduplicates by IP.
func Discover(ctx context.Context, timeout time.Duration) ([]Device, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	type result struct {
		devices []Device
		err     error
	}

	mdnsCh := make(chan result, 1)
	ssdpCh := make(chan result, 1)

	go func() {
		devs, err := DiscoverMDNS(ctx, timeout)
		mdnsCh <- result{devs, err}
	}()

	go func() {
		devs, err := DiscoverSSDP(ctx, timeout)
		ssdpCh <- result{devs, err}
	}()

	mdnsRes := <-mdnsCh
	ssdpRes := <-ssdpCh

	seen := make(map[string]struct{})
	var devices []Device

	// mDNS results first (typically more reliable).
	if mdnsRes.err == nil {
		for _, d := range mdnsRes.devices {
			if _, ok := seen[d.IP]; !ok {
				seen[d.IP] = struct{}{}
				devices = append(devices, d)
			}
		}
	}

	if ssdpRes.err == nil {
		for _, d := range ssdpRes.devices {
			if _, ok := seen[d.IP]; !ok {
				seen[d.IP] = struct{}{}
				devices = append(devices, d)
			}
		}
	}

	// Both failed → return first error.
	if mdnsRes.err != nil && ssdpRes.err != nil {
		return nil, fmt.Errorf("discovery failed: mdns: %w; ssdp: %v", mdnsRes.err, ssdpRes.err)
	}

	return devices, nil
}
