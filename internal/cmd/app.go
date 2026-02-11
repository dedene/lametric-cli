package cmd

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/dedene/lametric-cli/internal/api"
)

// AppCmd manages installed apps and built-in app controls.
type AppCmd struct {
	List      AppListCmd      `cmd:"" default:"1" help:"List installed apps"`
	Next      AppNextCmd      `cmd:"" help:"Switch to next app"`
	Prev      AppPrevCmd      `cmd:"" help:"Switch to previous app"`
	Activate  AppActivateCmd  `cmd:"" help:"Activate specific app widget"`
	Clock     AppClockCmd     `cmd:"" help:"Clock app controls"`
	Radio     AppRadioCmd     `cmd:"" help:"Radio app controls"`
	Timer     AppTimerCmd     `cmd:"" help:"Timer app controls"`
	Stopwatch AppStopwatchCmd `cmd:"" help:"Stopwatch controls"`
}

// AppListCmd lists installed apps.
type AppListCmd struct{}

func (c *AppListCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	var apps map[string]api.App
	if err := client.Get(context.Background(), "/api/v2/device/apps", &apps); err != nil {
		return fmt.Errorf("list apps: %w", err)
	}

	// Stable ordering by package name.
	keys := make([]string, 0, len(apps))
	for k := range apps {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	f := newFormatter(flags)
	headers := []string{"Package", "Vendor", "Version", "Widgets"}
	rows := make([][]string, 0, len(keys))
	for _, k := range keys {
		app := apps[k]
		rows = append(rows, []string{
			k,
			app.Vendor,
			app.Version,
			strconv.Itoa(len(app.Widgets)),
		})
	}

	return f.Output(apps, headers, rows)
}

// AppNextCmd switches to the next app.
type AppNextCmd struct{}

func (c *AppNextCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	if err := client.Put(context.Background(), "/api/v2/device/apps/next", nil, nil); err != nil {
		return fmt.Errorf("next app: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok"}, [][2]string{
		{"Status", "switched to next app"},
	})
}

// AppPrevCmd switches to the previous app.
type AppPrevCmd struct{}

func (c *AppPrevCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	if err := client.Put(context.Background(), "/api/v2/device/apps/prev", nil, nil); err != nil {
		return fmt.Errorf("prev app: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok"}, [][2]string{
		{"Status", "switched to previous app"},
	})
}

// AppActivateCmd activates a specific app widget.
type AppActivateCmd struct {
	Package  string `arg:"" help:"App package name (e.g. com.lametric.clock)"`
	WidgetID string `arg:"" help:"Widget ID to activate"`
}

func (c *AppActivateCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v2/device/apps/%s/widgets/%s/activate", c.Package, c.WidgetID)
	if err := client.Put(context.Background(), path, nil, nil); err != nil {
		return fmt.Errorf("activate app: %w", err)
	}

	f := newFormatter(flags)
	return f.OutputSingle(map[string]string{"status": "ok", "package": c.Package, "widget": c.WidgetID}, [][2]string{
		{"Status", "activated"},
		{"Package", c.Package},
		{"Widget", c.WidgetID},
	})
}

// appAction posts an action to a built-in app widget.
func appAction(flags *RootFlags, pkg, widgetID, actionID string, params map[string]any) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v2/device/apps/%s/widgets/%s/actions", pkg, widgetID)
	body := map[string]any{"id": actionID}
	if params != nil {
		body["params"] = params
	}

	if err := client.Post(context.Background(), path, body, nil); err != nil {
		return fmt.Errorf("%s: %w", actionID, err)
	}

	return nil
}

// firstWidgetID returns the first widget ID for a given package.
func firstWidgetID(flags *RootFlags, pkg string) (string, error) {
	client, _, err := resolveClient(flags)
	if err != nil {
		return "", err
	}

	var apps map[string]api.App
	if err := client.Get(context.Background(), "/api/v2/device/apps", &apps); err != nil {
		return "", fmt.Errorf("list apps: %w", err)
	}

	app, ok := apps[pkg]
	if !ok {
		return "", fmt.Errorf("app %q not found", pkg)
	}

	// Return first widget ID sorted by index.
	type wEntry struct {
		id    string
		index int
	}
	entries := make([]wEntry, 0, len(app.Widgets))
	for id, w := range app.Widgets {
		entries = append(entries, wEntry{id: id, index: w.Index})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].index < entries[j].index })

	if len(entries) == 0 {
		return "", fmt.Errorf("app %q has no widgets", pkg)
	}

	return entries[0].id, nil
}
