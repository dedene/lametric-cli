package cmd

const radioPkg = "com.lametric.radio"

// AppRadioCmd manages the radio app.
type AppRadioCmd struct {
	Play AppRadioPlayCmd `cmd:"" help:"Start playing radio"`
	Stop AppRadioStopCmd `cmd:"" help:"Stop radio"`
	Next AppRadioNextCmd `cmd:"" help:"Next station"`
	Prev AppRadioPrevCmd `cmd:"" help:"Previous station"`
}

// AppRadioPlayCmd starts playing the radio.
type AppRadioPlayCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppRadioPlayCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, radioPkg, wid, "radio.play", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "playing"}, [][2]string{
		{"Status", "playing"},
	})
}

func (c *AppRadioPlayCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, radioPkg)
}

// AppRadioStopCmd stops the radio.
type AppRadioStopCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppRadioStopCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, radioPkg, wid, "radio.stop", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "stopped"}, [][2]string{
		{"Status", "stopped"},
	})
}

func (c *AppRadioStopCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, radioPkg)
}

// AppRadioNextCmd switches to the next radio station.
type AppRadioNextCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppRadioNextCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, radioPkg, wid, "radio.next", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok"}, [][2]string{
		{"Status", "next station"},
	})
}

func (c *AppRadioNextCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, radioPkg)
}

// AppRadioPrevCmd switches to the previous radio station.
type AppRadioPrevCmd struct {
	Widget string `help:"Widget ID (auto-detected if omitted)" short:"w"`
}

func (c *AppRadioPrevCmd) Run(flags *RootFlags) error {
	wid, err := c.widgetID(flags)
	if err != nil {
		return err
	}

	if err := appAction(flags, radioPkg, wid, "radio.prev", nil); err != nil {
		return err
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok"}, [][2]string{
		{"Status", "previous station"},
	})
}

func (c *AppRadioPrevCmd) widgetID(flags *RootFlags) (string, error) {
	if c.Widget != "" {
		return c.Widget, nil
	}
	return firstWidgetID(flags, radioPkg)
}
