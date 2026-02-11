package discovery

import "testing"

func TestDevice_String(t *testing.T) {
	d := Device{Name: "LM1234", IP: "192.168.1.42", Port: 4343}
	got := d.String()
	want := "LM1234 (192.168.1.42:4343)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSanitizeMDNSName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"full mdns name", "LM1234._lametric-api._tcp.local.", "LM1234"},
		{"dotted", "LM1234.local.", "LM1234"},
		{"plain", "LM1234", "LM1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeMDNSName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeMDNSName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestModelFromName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"time model", "LM1234._lametric-api._tcp.local.", "TIME"},
		{"sky model", "LM_SKY_5678._lametric-api._tcp.local.", "SKY"},
		{"sky lowercase", "sky-device", "SKY"},
		{"default", "some-device", "TIME"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modelFromName(tt.input)
			if got != tt.want {
				t.Errorf("modelFromName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
