package cmd

import (
	"os"
	"sort"

	"github.com/dedene/lametric-cli/internal/output"
)

// IconAliases maps friendly names to LaMetric icon IDs.
var IconAliases = map[string]string{
	"rocket":     "i120",
	"checkmark":  "i2867",
	"cross":      "a2868",
	"heart":      "i339",
	"star":       "i2283",
	"fire":       "a2735",
	"warning":    "i555",
	"info":       "i65",
	"clock":      "i82",
	"calendar":   "i227",
	"mail":       "i65",
	"download":   "i120",
	"upload":     "i121",
	"github":     "i32320",
	"slack":      "i36682",
	"thumbsup":   "i14180",
	"thumbsdown": "i14182",
	"sun":        "i2282",
	"rain":       "i2284",
	"snow":       "i2289",
	"cloud":      "i2283",
	"music":      "i612",
	"beer":       "i915",
	"coffee":     "i16041",
	"pizza":      "i8779",
	"dollar":     "i73",
	"bitcoin":    "i12949",
	"eye":        "i14092",
	"lock":       "i242",
	"unlock":     "i1370",
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

// IconsCmd lists available icon aliases.
type IconsCmd struct{}

// Run prints the icon alias table.
func (c *IconsCmd) Run(flags *RootFlags) error {
	f := output.NewFormatter(os.Stdout, flags.JSON, flags.Plain, flags.NoColor)

	type iconEntry struct {
		Alias string `json:"alias"`
		ID    string `json:"id"`
	}

	names := make([]string, 0, len(IconAliases))
	for name := range IconAliases {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]iconEntry, 0, len(names))
	rows := make([][]string, 0, len(names))
	for _, name := range names {
		id := IconAliases[name]
		entries = append(entries, iconEntry{Alias: name, ID: id})
		rows = append(rows, []string{name, id})
	}

	return f.Output(entries, []string{"ALIAS", "ID"}, rows)
}
