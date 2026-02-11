package cmd

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
	"time"

	"golang.org/x/image/draw"

	"github.com/dedene/lametric-cli/internal/api"
)

// StreamCmd groups streaming subcommands.
type StreamCmd struct {
	Start  StreamStartCmd  `cmd:"" help:"Start streaming session"`
	Stop   StreamStopCmd   `cmd:"" help:"Stop streaming"`
	Status StreamStatusCmd `cmd:"" help:"Show streaming status"`
	Image  StreamImageCmd  `cmd:"" help:"Send single image"`
	Gif    StreamGifCmd    `cmd:"" help:"Send animated GIF"`
	Pipe   StreamPipeCmd   `cmd:"" help:"Pipe raw frames from stdin (for ffmpeg)"`
}

// StreamStartCmd starts a streaming session.
type StreamStartCmd struct{}

func (c *StreamStartCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	resp, err := client.StartStream(context.Background())
	if err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(resp, [][2]string{
		{"Session ID", resp.SessionID},
	})
}

// StreamStopCmd stops the active streaming session.
type StreamStopCmd struct{}

func (c *StreamStopCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	if err := client.StopStream(context.Background()); err != nil {
		return err
	}

	f := newFormatter(flags)

	type stopResult struct {
		Stopped bool `json:"stopped"`
	}

	return f.OutputSingle(stopResult{Stopped: true}, [][2]string{
		{"Status", "stopped"},
	})
}

// StreamStatusCmd shows current streaming status.
type StreamStatusCmd struct{}

func (c *StreamStatusCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	status, err := client.GetStreamStatus(context.Background())
	if err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(status, [][2]string{
		{"Active", fmt.Sprintf("%v", status.Active)},
		{"Session ID", status.SessionID},
	})
}

// StreamImageCmd sends a single image to the device.
type StreamImageCmd struct {
	File string `arg:"" help:"Image file (PNG, JPG)"`
}

func (c *StreamImageCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	// Get device display dimensions.
	width, height, err := getDeviceDimensions(client)
	if err != nil {
		return err
	}

	// Start streaming session.
	resp, err := client.StartStream(context.Background())
	if err != nil {
		return err
	}
	defer func() { _ = client.StopStream(context.Background()) }()

	session, err := api.NewStreamSession(extractIP(client), resp.SessionID, width, height)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close() }()

	// Load and scale image.
	imgFile, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer func() { _ = imgFile.Close() }()

	src, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	scaled := scaleImage(src, width, height)

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return fmt.Errorf("encode PNG: %w", err)
	}

	if err := session.SendFrame(api.EncodingPNG, buf.Bytes()); err != nil {
		return err
	}

	f := newFormatter(flags)

	type imgResult struct {
		Sent bool `json:"sent"`
	}

	return f.OutputSingle(imgResult{Sent: true}, [][2]string{
		{"Status", "sent"},
		{"Size", fmt.Sprintf("%dx%d", width, height)},
	})
}

// StreamGifCmd sends an animated GIF frame-by-frame.
type StreamGifCmd struct {
	File string `arg:"" help:"GIF file"`
	Loop bool   `help:"Loop the GIF continuously" short:"l"`
}

func (c *StreamGifCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	width, height, err := getDeviceDimensions(client)
	if err != nil {
		return err
	}

	resp, err := client.StartStream(context.Background())
	if err != nil {
		return err
	}
	defer func() { _ = client.StopStream(context.Background()) }()

	session, err := api.NewStreamSession(extractIP(client), resp.SessionID, width, height)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close() }()

	gifFile, err := os.Open(c.File)
	if err != nil {
		return fmt.Errorf("open GIF: %w", err)
	}
	defer func() { _ = gifFile.Close() }()

	g, err := gif.DecodeAll(gifFile)
	if err != nil {
		return fmt.Errorf("decode GIF: %w", err)
	}

	for {
		for i, frame := range g.Image {
			scaled := scaleImage(frame, width, height)

			var buf bytes.Buffer
			if err := png.Encode(&buf, scaled); err != nil {
				return fmt.Errorf("encode frame %d: %w", i, err)
			}

			if err := session.SendFrame(api.EncodingPNG, buf.Bytes()); err != nil {
				return err
			}

			delay := 100 * time.Millisecond // default 100ms
			if i < len(g.Delay) && g.Delay[i] > 0 {
				delay = time.Duration(g.Delay[i]) * 10 * time.Millisecond
			}
			time.Sleep(delay)
		}

		if !c.Loop {
			break
		}
	}

	return nil
}

// StreamPipeCmd reads raw RGB888 frames from stdin and streams them.
type StreamPipeCmd struct {
	Width  int `help:"Frame width (default: device width)" short:"w"`
	Height int `help:"Frame height (default: device height)" short:"h"`
	FPS    int `help:"Target frames per second" default:"30"`
}

func (c *StreamPipeCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	width, height := c.Width, c.Height
	if width == 0 || height == 0 {
		w, h, dimErr := getDeviceDimensions(client)
		if dimErr != nil {
			return dimErr
		}
		if width == 0 {
			width = w
		}
		if height == 0 {
			height = h
		}
	}

	resp, err := client.StartStream(context.Background())
	if err != nil {
		return err
	}
	defer func() { _ = client.StopStream(context.Background()) }()

	session, err := api.NewStreamSession(extractIP(client), resp.SessionID, width, height)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close() }()

	frameSize := width * height * 3 // RGB888
	frameBuf := make([]byte, frameSize)
	interval := time.Second / time.Duration(c.FPS)

	for {
		_, err := io.ReadFull(os.Stdin, frameBuf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			return fmt.Errorf("read frame: %w", err)
		}

		start := time.Now()

		if err := session.SendFrame(api.EncodingRAW, frameBuf); err != nil {
			return err
		}

		elapsed := time.Since(start)
		if elapsed < interval {
			time.Sleep(interval - elapsed)
		}
	}
}

// getDeviceDimensions fetches display width/height from the device.
func getDeviceDimensions(client *api.Client) (int, int, error) {
	var device api.Device
	if err := client.Get(context.Background(), "/api/v2/device", &device); err != nil {
		return 0, 0, fmt.Errorf("get device info: %w", err)
	}

	w, h := device.Display.Width, device.Display.Height
	if w == 0 {
		w = 37 // TIME default
	}
	if h == 0 {
		h = 8
	}
	return w, h, nil
}

// extractIP returns the device IP from the client's base URL.
func extractIP(client *api.Client) string {
	return client.IP()
}

// scaleImage resizes an image to the target dimensions using Catmull-Rom.
func scaleImage(src image.Image, width, height int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
