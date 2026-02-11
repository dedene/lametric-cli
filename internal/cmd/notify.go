package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dedene/lametric-cli/internal/api"
)

// NotifyCmd sends a notification to a LaMetric device.
type NotifyCmd struct {
	Text     string `arg:"" optional:"" help:"Notification text (or use stdin)"`
	Icon     string `help:"Icon ID or alias (e.g., rocket, i1234)" short:"i"`
	Sound    string `help:"Sound ID or alias (e.g., positive1)" short:"s"`
	Priority string `help:"Priority: info, warning, critical" default:"info" short:"p" enum:"info,warning,critical"`
	Cycles   int    `help:"Repeat count (0=loop)" default:"1"`
	Lifetime int    `help:"TTL in milliseconds" default:"120000"`
	Goal     string `help:"Goal frame: current/max (e.g., 50/100)"`
	Chart    string `help:"Chart frame: comma-separated values (e.g., 1,2,3,4,5)"`
	Wait     bool   `help:"Block until notification dismissed"`
}

// Run executes the notify command.
func (c *NotifyCmd) Run(flags *RootFlags) error {
	text := c.Text
	if text == "" {
		text = readStdin()
	}

	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	frames, err := c.buildFrames(text)
	if err != nil {
		return err
	}

	req := api.NotificationRequest{
		Priority: c.Priority,
		Lifetime: c.Lifetime,
		Model: api.NotificationModel{
			Frames: frames,
			Sound:  ResolveSound(c.Sound),
			Cycles: c.Cycles,
		},
	}

	var notif api.Notification
	if err := client.Post(context.Background(), "/api/v2/device/notifications", req, &notif); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	if c.Wait {
		if err := c.waitForDismissal(client, notif.ID); err != nil {
			return err
		}
	}

	f := newFormatter(flags)

	type notifyResult struct {
		ID      string `json:"id"`
		Success bool   `json:"success"`
	}

	return f.OutputSingle(notifyResult{ID: notif.ID, Success: true}, [][2]string{
		{"ID", notif.ID},
		{"Status", "sent"},
	})
}

func (c *NotifyCmd) buildFrames(text string) ([]api.Frame, error) {
	var frames []api.Frame

	// Text frame (always first if text is present).
	if text != "" {
		frames = append(frames, api.Frame{
			Icon: ResolveIcon(c.Icon),
			Text: text,
		})
	}

	// Goal frame.
	if c.Goal != "" {
		gd, err := parseGoal(c.Goal)
		if err != nil {
			return nil, err
		}

		frame := api.Frame{GoalData: gd, Icon: ResolveIcon(c.Icon)}
		frames = append(frames, frame)
	}

	// Chart frame.
	if c.Chart != "" {
		data, err := parseChart(c.Chart)
		if err != nil {
			return nil, err
		}

		frames = append(frames, api.Frame{ChartData: data})
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no notification content: provide text, --goal, or --chart")
	}

	return frames, nil
}

func (c *NotifyCmd) waitForDismissal(client *api.Client, id string) error {
	path := "/api/v2/device/notifications/" + id
	for {
		var notif api.Notification
		err := client.Get(context.Background(), path, &notif)
		if err != nil {
			// Notification gone = dismissed.
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func readStdin() string {
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice != 0 {
		return ""
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.TrimSpace(strings.Join(lines, " "))
}

// parseGoal parses "current/max" into GoalData.
func parseGoal(s string) (*api.GoalData, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid goal format %q: expected current/max (e.g., 50/100)", s)
	}

	current, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid goal current %q: %w", parts[0], err)
	}

	max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid goal max %q: %w", parts[1], err)
	}

	return &api.GoalData{Start: 0, Current: current, End: max}, nil
}

// parseChart parses comma-separated ints.
func parseChart(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	data := make([]int, 0, len(parts))

	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid chart value %q: %w", p, err)
		}

		data = append(data, v)
	}

	return data, nil
}
