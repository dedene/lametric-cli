package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dedene/lametric-cli/internal/api"
)

// AudioCmd manages audio settings.
type AudioCmd struct {
	Get    AudioGetCmd    `cmd:"" default:"1" help:"Show audio settings"`
	Volume AudioVolumeCmd `cmd:"" help:"Set volume (0-100)"`
}

// AudioGetCmd shows current audio settings.
type AudioGetCmd struct{}

func (c *AudioGetCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var audio api.Audio
	if err := client.Get(context.Background(), "/api/v2/device/audio", &audio); err != nil {
		return fmt.Errorf("get audio: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(audio, [][2]string{
		{"Volume", strconv.Itoa(audio.Volume)},
	})
}

// AudioVolumeCmd sets the volume level.
type AudioVolumeCmd struct {
	Level int `arg:"" help:"Volume level 0-100"`
}

func (c *AudioVolumeCmd) Run(flags *RootFlags) error {
	if c.Level < 0 || c.Level > 100 {
		return fmt.Errorf("volume must be 0-100, got %d", c.Level)
	}

	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	body := api.AudioUpdate{Volume: &c.Level}

	var audio api.Audio
	if err := client.Put(context.Background(), "/api/v2/device/audio", body, &audio); err != nil {
		return fmt.Errorf("set volume: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(audio, [][2]string{
		{"Volume", strconv.Itoa(audio.Volume)},
	})
}
