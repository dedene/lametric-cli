package discovery

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/koron/go-ssdp"
)

var ssdpSearchTargets = []string{
	"urn:schemas-upnp-org:device:LaMetric:1",
	"urn:schemas-upnp-org:device:LaMetric:2",
}

// DiscoverSSDP finds LaMetric devices via UPnP/SSDP.
func DiscoverSSDP(ctx context.Context, timeout time.Duration) ([]Device, error) {
	var all []Device

	for _, target := range ssdpSearchTargets {
		if ctx.Err() != nil {
			break
		}

		services, err := ssdp.Search(target, int(timeout.Seconds()), "")
		if err != nil {
			continue
		}

		for _, svc := range services {
			d, err := deviceFromLocation(ctx, svc.Location)
			if err != nil {
				continue
			}
			all = append(all, d)
		}
	}

	return all, nil
}

// ssdpRoot is the minimal UPnP device description XML.
type ssdpRoot struct {
	XMLName xml.Name   `xml:"root"`
	Device  ssdpDevice `xml:"device"`
}

type ssdpDevice struct {
	FriendlyName string `xml:"friendlyName"`
	ModelName    string `xml:"modelName"`
	Manufacturer string `xml:"manufacturer"`
}

// deviceFromLocation fetches a UPnP device description and parses it.
func deviceFromLocation(ctx context.Context, location string) (Device, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, location, nil)
	if err != nil {
		return Device{}, fmt.Errorf("ssdp request: %w", err)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Device{}, fmt.Errorf("ssdp fetch: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return Device{}, fmt.Errorf("ssdp read: %w", err)
	}

	var root ssdpRoot
	if err := xml.Unmarshal(body, &root); err != nil {
		return Device{}, fmt.Errorf("ssdp xml: %w", err)
	}

	// Validate this is actually a LaMetric device
	manufacturer := strings.ToLower(root.Device.Manufacturer)
	modelName := strings.ToLower(root.Device.ModelName)
	if !strings.Contains(manufacturer, "lametric") &&
		!strings.Contains(modelName, "lametric") {
		return Device{}, fmt.Errorf("not a LaMetric device: %s", root.Device.FriendlyName)
	}

	host, port, _ := net.SplitHostPort(req.URL.Host)
	if host == "" {
		host = req.URL.Hostname()
	}

	p := 443
	if port != "" {
		if v, err := strconv.Atoi(port); err == nil {
			p = v
		}
	}

	model := "TIME"
	if strings.Contains(strings.ToLower(root.Device.ModelName), "sky") {
		model = "SKY"
	}

	return Device{
		Name:  root.Device.FriendlyName,
		IP:    host,
		Port:  p,
		Model: model,
	}, nil
}
