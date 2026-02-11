package cmd

const clockPkg = "com.lametric.clock"

// AppClockCmd manages the clock app.
type AppClockCmd struct {
	Alarm AppClockAlarmCmd `cmd:"" help:"Toggle clock alarm"`
	Show  AppClockShowCmd  `cmd:"" help:"Show clock face"`
}

// AppClockAlarmCmd toggles the clock alarm.
type AppClockAlarmCmd struct {
	Enabled bool   `arg:"" help:"Enable (true) or disable (false) alarm"`
	Widget  string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppClockAlarmCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	actionID := "clock.alarm"
	params := map[string]any{"enabled": c.Enabled}

	if err := appAction(flags, clockPkg, wid, actionID, params); err != nil {
		return err
	}

	f := newFormatter(flags)
	status := "disabled"
	if c.Enabled {
		status = "enabled"
	}
	return f.OutputSingle(map[string]string{"alarm": status}, [][2]string{
		{"Alarm", status},
	})
}

func (c *AppClockAlarmCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, clockPkg)
}

// AppClockShowCmd activates the clock face.
type AppClockShowCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppClockShowCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, clockPkg, wid, "clock.clockface", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok"}, [][2]string{
		{"Status", "clock face shown"},
	})
}

func (c *AppClockShowCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, clockPkg)
}
