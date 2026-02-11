package cmd

import (
	"fmt"
	"time"
)

const timerPkg = "com.lametric.countdown"

// AppTimerCmd manages the countdown timer app.
type AppTimerCmd struct {
	Set   AppTimerSetCmd   `cmd:"" help:"Set timer duration (e.g. 5m, 1h30m, 90s)"`
	Start AppTimerStartCmd `cmd:"" help:"Start timer"`
	Pause AppTimerPauseCmd `cmd:"" help:"Pause timer"`
	Reset AppTimerResetCmd `cmd:"" help:"Reset timer"`
}

// AppTimerSetCmd sets the timer duration and starts it.
type AppTimerSetCmd struct {
	Duration string `arg:"" help:"Duration (e.g. 5m, 1h30m, 90s)"`
	Widget   string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppTimerSetCmd) Run(flags *RootFlags) error {
	d, err := time.ParseDuration(c.Duration)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", c.Duration, err)
	}

	secs := int(d.Seconds())
	if secs <= 0 {
		return fmt.Errorf("duration must be positive, got %s", c.Duration)
	}

	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	params := map[string]any{"duration": secs, "start_now": true}
	if err := appAction(flags, timerPkg, wid, "countdown.configure", params); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]any{"duration": secs, "status": "started"}, [][2]string{
		{"Duration", d.String()},
		{"Status", "started"},
	})
}

func (c *AppTimerSetCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, timerPkg)
}

// AppTimerStartCmd starts the timer.
type AppTimerStartCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppTimerStartCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, timerPkg, wid, "countdown.start", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "started"}, [][2]string{
		{"Status", "started"},
	})
}

func (c *AppTimerStartCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, timerPkg)
}

// AppTimerPauseCmd pauses the timer.
type AppTimerPauseCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppTimerPauseCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, timerPkg, wid, "countdown.pause", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "paused"}, [][2]string{
		{"Status", "paused"},
	})
}

func (c *AppTimerPauseCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, timerPkg)
}

// AppTimerResetCmd resets the timer.
type AppTimerResetCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppTimerResetCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, timerPkg, wid, "countdown.reset", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "reset"}, [][2]string{
		{"Status", "reset"},
	})
}

func (c *AppTimerResetCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, timerPkg)
}
