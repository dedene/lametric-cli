package cmd

const stopwatchPkg = "com.lametric.stopwatch"

// AppStopwatchCmd manages the stopwatch app.
type AppStopwatchCmd struct {
	Start AppStopwatchStartCmd `cmd:"" help:"Start stopwatch"`
	Pause AppStopwatchPauseCmd `cmd:"" help:"Pause stopwatch"`
	Reset AppStopwatchResetCmd `cmd:"" help:"Reset stopwatch"`
}

// AppStopwatchStartCmd starts the stopwatch.
type AppStopwatchStartCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppStopwatchStartCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, stopwatchPkg, wid, "stopwatch.start", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "started"}, [][2]string{
		{"Status", "started"},
	})
}

func (c *AppStopwatchStartCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, stopwatchPkg)
}

// AppStopwatchPauseCmd pauses the stopwatch.
type AppStopwatchPauseCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppStopwatchPauseCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, stopwatchPkg, wid, "stopwatch.pause", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "paused"}, [][2]string{
		{"Status", "paused"},
	})
}

func (c *AppStopwatchPauseCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, stopwatchPkg)
}

// AppStopwatchResetCmd resets the stopwatch.
type AppStopwatchResetCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppStopwatchResetCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, stopwatchPkg, wid, "stopwatch.reset", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "reset"}, [][2]string{
		{"Status", "reset"},
	})
}

func (c *AppStopwatchResetCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, stopwatchPkg)
}
