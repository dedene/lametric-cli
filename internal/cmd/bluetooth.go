package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dedene/lametric-cli/internal/api"
)

// BluetoothCmd manages bluetooth settings.
type BluetoothCmd struct {
	Get  BluetoothGetCmd  `cmd:"" default:"1" help:"Show bluetooth status"`
	On   BluetoothOnCmd   `cmd:"" help:"Enable bluetooth"`
	Off  BluetoothOffCmd  `cmd:"" help:"Disable bluetooth"`
	Name BluetoothNameCmd `cmd:"" help:"Set bluetooth name"`
}

// BluetoothGetCmd shows current bluetooth status.
type BluetoothGetCmd struct{}

func (c *BluetoothGetCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var bt api.Bluetooth
	if err := client.Get(context.Background(), "/api/v2/device/bluetooth", &bt); err != nil {
		return fmt.Errorf("get bluetooth: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(bt, [][2]string{
		{"Active", strconv.FormatBool(bt.Active)},
		{"Available", strconv.FormatBool(bt.Available)},
		{"Name", bt.Name},
		{"MAC", bt.MAC},
		{"Pairable", strconv.FormatBool(bt.Pairable)},
		{"Discoverable", strconv.FormatBool(bt.Discoverable)},
	})
}

// BluetoothOnCmd enables bluetooth.
type BluetoothOnCmd struct{}

func (c *BluetoothOnCmd) Run(flags *RootFlags) error {
	return setBluetoothActive(flags, true)
}

// BluetoothOffCmd disables bluetooth.
type BluetoothOffCmd struct{}

func (c *BluetoothOffCmd) Run(flags *RootFlags) error {
	return setBluetoothActive(flags, false)
}

func setBluetoothActive(flags *RootFlags, active bool) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	body := api.BluetoothUpdate{Active: &active}

	var bt api.Bluetooth
	if err := client.Put(context.Background(), "/api/v2/device/bluetooth", body, &bt); err != nil {
		return fmt.Errorf("set bluetooth: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(bt, [][2]string{
		{"Active", strconv.FormatBool(bt.Active)},
		{"Name", bt.Name},
	})
}

// BluetoothNameCmd sets the bluetooth device name.
type BluetoothNameCmd struct {
	Name string `arg:"" help:"Bluetooth device name"`
}

func (c *BluetoothNameCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	body := api.BluetoothUpdate{Name: &c.Name}

	var bt api.Bluetooth
	if err := client.Put(context.Background(), "/api/v2/device/bluetooth", body, &bt); err != nil {
		return fmt.Errorf("set bluetooth name: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(bt, [][2]string{
		{"Active", strconv.FormatBool(bt.Active)},
		{"Name", bt.Name},
	})
}
