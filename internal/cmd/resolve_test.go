package cmd

import (
	"testing"

	"github.com/dedene/lametric-cli/internal/api"
)

func TestResolveIcon(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"alias rocket", "rocket", "a26304"},
		{"alias heart", "heart", "a230"},
		{"icon id i-prefix", "i9999", "i9999"},
		{"icon id a-prefix", "a1234", "a1234"},
		{"unknown passthrough", "custom", "custom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveIcon(tt.input)
			if got != tt.want {
				t.Errorf("ResolveIcon(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestResolveSound(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *api.Sound
	}{
		{"empty", "", nil},
		{"alias positive1", "positive1", &api.Sound{Category: "notifications", ID: "positive1"}},
		{"alias alarm1", "alarm1", &api.Sound{Category: "alarms", ID: "alarm1"}},
		{"unknown defaults to notifications", "custom_sound", &api.Sound{Category: "notifications", ID: "custom_sound"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveSound(tt.input)
			if tt.want == nil {
				if got != nil {
					t.Errorf("ResolveSound(%q) = %+v, want nil", tt.input, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("ResolveSound(%q) = nil, want %+v", tt.input, tt.want)
			}
			if got.Category != tt.want.Category || got.ID != tt.want.ID {
				t.Errorf("ResolveSound(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}
