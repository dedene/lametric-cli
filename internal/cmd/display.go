package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dedene/lametric-cli/internal/api"
)

// DisplayCmd manages display settings.
type DisplayCmd struct {
	Get        DisplayGetCmd        `cmd:"" default:"1" help:"Show display settings"`
	Brightness DisplayBrightnessCmd `cmd:"" help:"Set brightness (0-100)"`
	Mode       DisplayModeCmd       `cmd:"" help:"Set brightness mode (auto/manual)"`
}

// DisplayGetCmd shows current display settings.
type DisplayGetCmd struct{}

func (c *DisplayGetCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var display api.Display
	if err := client.Get(context.Background(), "/api/v2/device/display", &display); err != nil {
		return fmt.Errorf("get display: %w", err)
	}

	f := newFormatter(flags)

	screensaver := "off"
	if display.Screensaver != nil && display.Screensaver.Enabled {
		screensaver = "on"
	}

	return f.OutputSingle(display, [][2]string{
		{"Brightness", strconv.Itoa(display.Brightness)},
		{"Mode", display.BrightnessMode},
		{"Width", strconv.Itoa(display.Width)},
		{"Height", strconv.Itoa(display.Height)},
		{"Type", display.Type},
		{"Screensaver", screensaver},
	})
}

// DisplayBrightnessCmd sets the display brightness.
type DisplayBrightnessCmd struct {
	Level int `arg:"" help:"Brightness level 0-100"`
}

func (c *DisplayBrightnessCmd) Run(flags *RootFlags) error {
	if c.Level < 0 || c.Level > 100 {
		return fmt.Errorf("brightness must be 0-100, got %d", c.Level)
	}

	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	body := api.DisplayUpdate{Brightness: &c.Level}

	var display api.Display
	if err := client.Put(context.Background(), "/api/v2/device/display", body, &display); err != nil {
		return fmt.Errorf("set brightness: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(display, [][2]string{
		{"Brightness", strconv.Itoa(display.Brightness)},
		{"Mode", display.BrightnessMode},
	})
}

// DisplayModeCmd sets the brightness mode.
type DisplayModeCmd struct {
	Mode string `arg:"" help:"Brightness mode: auto or manual" enum:"auto,manual"`
}

func (c *DisplayModeCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	body := api.DisplayUpdate{BrightnessMode: &c.Mode}

	var display api.Display
	if err := client.Put(context.Background(), "/api/v2/device/display", body, &display); err != nil {
		return fmt.Errorf("set brightness mode: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(display, [][2]string{
		{"Brightness", strconv.Itoa(display.Brightness)},
		{"Mode", display.BrightnessMode},
	})
}
