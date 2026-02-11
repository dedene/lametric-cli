package errfmt

import (
	"errors"
	"strings"
	"testing"
)

func TestFormat_Nil(t *testing.T) {
	if got := Format(nil); got != "" {
		t.Errorf("Format(nil) = %q, want empty", got)
	}
}

func TestFormat_KnownPatterns(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		contain string
	}{
		{"no api key", errors.New("no api key configured"), "lametric setup"},
		{"unauthorized", errors.New("unauthorized access"), "Authentication failed"},
		{"connection refused", errors.New("connection refused"), "powered on"},
		{"timeout", errors.New("request timeout"), "timed out"},
		{"not found", errors.New("resource not found"), "not found"},
		{"tls", errors.New("tls handshake error"), "TLS error"},
		{"certificate", errors.New("certificate verify failed"), "Certificate error"},
		{"no such host", errors.New("no such host"), "resolve device"},
		{"forbidden", errors.New("forbidden"), "Access denied"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.err)
			if !strings.Contains(got, tt.contain) {
				t.Errorf("Format(%q) = %q, want to contain %q", tt.err, got, tt.contain)
			}
		})
	}
}

func TestFormat_UnknownError(t *testing.T) {
	err := errors.New("something completely unexpected")
	got := Format(err)
	if got != "something completely unexpected" {
		t.Errorf("got %q, want original message", got)
	}
}

func TestWrap_Nil(t *testing.T) {
	if got := Wrap(nil, "test"); got != nil {
		t.Errorf("Wrap(nil) = %v, want nil", got)
	}
}

func TestWrap(t *testing.T) {
	err := Wrap(errors.New("boom"), "send notification")
	if !strings.Contains(err.Error(), "send notification") || !strings.Contains(err.Error(), "boom") {
		t.Errorf("Wrap() = %q", err)
	}
}

func TestHint_Nil(t *testing.T) {
	if got := Hint(nil, "hint"); got != "" {
		t.Errorf("Hint(nil) = %q, want empty", got)
	}
}

func TestHint(t *testing.T) {
	err := errors.New("connection refused")
	got := Hint(err, "check firewall")
	if !strings.Contains(got, "hint:") || !strings.Contains(got, "check firewall") {
		t.Errorf("Hint() = %q", got)
	}
}
