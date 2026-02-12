package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dedene/lametric-cli/internal/api"
	"github.com/dedene/lametric-cli/internal/output"
)

// IconAliases maps friendly names to LaMetric icon IDs.
// Use "lametric icons search <query>" to find more icons.
var IconAliases = map[string]string{
	"checkmark":  "i2867",
	"cross":      "a2868",
	"heart":      "a230",
	"star":       "i635",
	"fire":       "a2735",
	"warning":    "i555",
	"info":       "i65",
	"clock":      "a82",
	"calendar":   "i66",
	"mail":       "i43",
	"email":      "i20907",
	"download":   "i120",
	"upload":     "i121",
	"rocket":     "a26304",
	"github":     "i32320",
	"slack":      "i36682",
	"thumbsup":   "i14180",
	"thumbsdown": "i14182",
	"sun":        "i2282",
	"rain":       "i2284",
	"snow":       "i2289",
	"cloud":      "i2283",
	"music":      "i612",
	"beer":       "i3253",
	"coffee":     "i16041",
	"pizza":      "i8779",
	"dollar":     "i73",
	"bitcoin":    "i12949",
	"eye":        "i14092",
	"lock":       "i242",
	"unlock":     "i1370",
	"bell":       "a1370",
	"phone":      "i124",
	"home":       "i96",
	"error":      "i555",
	"success":    "i2867",
	"fail":       "a2868",
}

// ResolveIcon returns the icon ID for a name or alias.
// If the input already looks like an icon ID (starts with "i" or "a"),
// it is returned as-is.
func ResolveIcon(nameOrID string) string {
	if nameOrID == "" {
		return ""
	}
	if nameOrID[0] == 'i' || nameOrID[0] == 'a' {
		return nameOrID
	}
	if id, ok := IconAliases[nameOrID]; ok {
		return id
	}
	return nameOrID
}

// IconsCmd manages icon search.
type IconsCmd struct {
	Popular IconsPopularCmd `cmd:"" default:"1" help:"Show popular icons"`
	Search  IconsSearchCmd  `cmd:"" help:"Search icons by name from LaMetric cloud"`
}

// IconsPopularCmd shows popular icons from the cloud.
type IconsPopularCmd struct {
	Limit int `help:"Maximum results to show" default:"30" short:"n"`
}

// Run fetches and displays popular icons.
func (c *IconsPopularCmd) Run(flags *RootFlags) error {
	f := output.NewFormatter(os.Stdout, flags.JSON, flags.Plain, flags.NoColor)

	spinner := output.NewSpinner("Fetching popular icons...")
	spinner.Start()

	icons, err := api.GetPopularIcons(context.Background(), c.Limit)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("fetch icons: %w", err)
	}

	type iconResult struct {
		Code  string `json:"code"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}

	results := make([]iconResult, len(icons))
	rows := make([][]string, len(icons))
	for i, icon := range icons {
		results[i] = iconResult{Code: icon.Code, Title: icon.Title, Type: icon.Type}
		rows[i] = []string{icon.Code, icon.Title, icon.Type}
	}

	return f.Output(results, []string{"CODE", "TITLE", "TYPE"}, rows)
}

// IconsSearchCmd searches icons in the LaMetric cloud library.
type IconsSearchCmd struct {
	Query string `arg:"" help:"Search query (e.g., 'rocket', 'heart')"`
	Limit int    `help:"Maximum results to show" default:"20" short:"n"`
}

// Run searches for icons matching the query.
func (c *IconsSearchCmd) Run(flags *RootFlags) error {
	f := output.NewFormatter(os.Stdout, flags.JSON, flags.Plain, flags.NoColor)

	spinner := output.NewSpinner("Searching icons...")
	spinner.Start()

	icons, err := api.SearchIcons(context.Background(), c.Query, c.Limit)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("search icons: %w", err)
	}

	if len(icons) == 0 {
		fmt.Printf("No icons found matching %q\n", c.Query)
		return nil
	}

	type iconResult struct {
		Code  string `json:"code"`
		Title string `json:"title"`
		Type  string `json:"type"`
	}

	results := make([]iconResult, len(icons))
	rows := make([][]string, len(icons))
	for i, icon := range icons {
		results[i] = iconResult{Code: icon.Code, Title: icon.Title, Type: icon.Type}
		rows[i] = []string{icon.Code, icon.Title, icon.Type}
	}

	return f.Output(results, []string{"CODE", "TITLE", "TYPE"}, rows)
}
