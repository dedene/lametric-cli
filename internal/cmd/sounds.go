package cmd

import (
	"os"
	"sort"

	"github.com/dedene/lametric-cli/internal/api"
	"github.com/dedene/lametric-cli/internal/output"
)

// SoundAliases maps friendly names to Sound definitions.
var SoundAliases = map[string]api.Sound{
	// Notifications
	"notification":  {Category: "notifications", ID: "notification"},
	"notification2": {Category: "notifications", ID: "notification2"},
	"notification3": {Category: "notifications", ID: "notification3"},
	"notification4": {Category: "notifications", ID: "notification4"},
	"positive1":     {Category: "notifications", ID: "positive1"},
	"positive2":     {Category: "notifications", ID: "positive2"},
	"positive3":     {Category: "notifications", ID: "positive3"},
	"positive4":     {Category: "notifications", ID: "positive4"},
	"positive5":     {Category: "notifications", ID: "positive5"},
	"positive6":     {Category: "notifications", ID: "positive6"},
	"negative1":     {Category: "notifications", ID: "negative1"},
	"negative2":     {Category: "notifications", ID: "negative2"},
	"negative3":     {Category: "notifications", ID: "negative3"},
	"negative4":     {Category: "notifications", ID: "negative4"},
	"negative5":     {Category: "notifications", ID: "negative5"},
	"win":           {Category: "notifications", ID: "win"},
	"win2":          {Category: "notifications", ID: "win2"},
	"cash":          {Category: "notifications", ID: "cash"},
	"cat":           {Category: "notifications", ID: "cat"},
	"dog":           {Category: "notifications", ID: "dog"},
	"bicycle":       {Category: "notifications", ID: "bicycle"},
	"energy":        {Category: "notifications", ID: "energy"},
	"knock":         {Category: "notifications", ID: "knock-knock"},
	"letter":        {Category: "notifications", ID: "letter_email"},
	"lose1":         {Category: "notifications", ID: "lose1"},
	"lose2":         {Category: "notifications", ID: "lose2"},
	"statistic":     {Category: "notifications", ID: "statistic"},
	"wind":          {Category: "notifications", ID: "wind"},
	"wind2":         {Category: "notifications", ID: "wind_short"},

	// Alarms
	"alarm1":  {Category: "alarms", ID: "alarm1"},
	"alarm2":  {Category: "alarms", ID: "alarm2"},
	"alarm3":  {Category: "alarms", ID: "alarm3"},
	"alarm4":  {Category: "alarms", ID: "alarm4"},
	"alarm5":  {Category: "alarms", ID: "alarm5"},
	"alarm6":  {Category: "alarms", ID: "alarm6"},
	"alarm7":  {Category: "alarms", ID: "alarm7"},
	"alarm8":  {Category: "alarms", ID: "alarm8"},
	"alarm9":  {Category: "alarms", ID: "alarm9"},
	"alarm10": {Category: "alarms", ID: "alarm10"},
	"alarm11": {Category: "alarms", ID: "alarm11"},
	"alarm12": {Category: "alarms", ID: "alarm12"},
	"alarm13": {Category: "alarms", ID: "alarm13"},
}

// ResolveSound returns the Sound for a name or alias.
// If the name is not found, it's treated as a notification sound ID.
func ResolveSound(name string) *api.Sound {
	if name == "" {
		return nil
	}
	if s, ok := SoundAliases[name]; ok {
		return &s
	}
	return &api.Sound{Category: "notifications", ID: name}
}

// SoundsCmd lists available sound aliases.
type SoundsCmd struct{}

// Run prints the sound alias table.
func (c *SoundsCmd) Run(flags *RootFlags) error {
	f := output.NewFormatter(os.Stdout, flags.JSON, flags.Plain, flags.NoColor)

	type soundEntry struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		ID       string `json:"id"`
	}

	names := make([]string, 0, len(SoundAliases))
	for name := range SoundAliases {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]soundEntry, 0, len(names))
	rows := make([][]string, 0, len(names))
	for _, name := range names {
		s := SoundAliases[name]
		entries = append(entries, soundEntry{Name: name, Category: s.Category, ID: s.ID})
		rows = append(rows, []string{name, s.Category, s.ID})
	}

	return f.Output(entries, []string{"NAME", "CATEGORY", "ID"}, rows)
}
