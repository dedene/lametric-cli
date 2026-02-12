package discovery

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

const mdnsService = "_lametric-api._tcp"

// DiscoverMDNS finds LaMetric devices via mDNS.
func DiscoverMDNS(ctx context.Context, timeout time.Duration) ([]Device, error) {
	// Suppress mdns library INFO logging
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	entriesCh := make(chan *mdns.ServiceEntry, 16)
	var devices []Device

	done := make(chan struct{})
	go func() {
		defer close(done)
		for entry := range entriesCh {
			// Determine IP address (prefer IPv4).
			var ip string
			switch {
			case entry.AddrV4 != nil:
				ip = entry.AddrV4.String()
			case entry.AddrV6 != nil:
				ip = entry.AddrV6.String()
			default:
				continue // Skip entries without any IP
			}

			d := Device{
				Name:  sanitizeMDNSName(entry.Name),
				IP:    ip,
				Port:  entry.Port,
				Model: modelFromName(entry.Name),
			}
			devices = append(devices, d)
		}
	}()

	params := mdns.DefaultParams(mdnsService)
	params.Entries = entriesCh
	params.Timeout = timeout
	params.DisableIPv6 = true

	if ctx.Err() != nil {
		close(entriesCh)
		<-done
		return nil, ctx.Err()
	}

	err := mdns.Query(params)
	close(entriesCh)
	<-done

	if err != nil {
		return nil, fmt.Errorf("mdns query: %w", err)
	}

	return devices, nil
}

// sanitizeMDNSName extracts the instance name from a full mDNS name.
func sanitizeMDNSName(name string) string {
	// mDNS names look like "LM1234._lametric-api._tcp.local."
	if idx := strings.Index(name, "."+mdnsService); idx > 0 {
		return name[:idx]
	}
	if idx := strings.Index(name, "."); idx > 0 {
		return name[:idx]
	}
	return name
}

// modelFromName guesses the model from the mDNS instance name.
func modelFromName(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "sky"):
		return "SKY"
	default:
		return "TIME"
	}
}
