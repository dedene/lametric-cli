package cmd

import (
	"errors"
	"strings"
	"testing"
)

func TestExitError_Error(t *testing.T) {
	err := &ExitError{Code: 1, Err: errors.New("failed")}
	if err.Error() != "failed" {
		t.Errorf("Error() = %q, want %q", err.Error(), "failed")
	}
}

func TestExitError_Error_Nil(t *testing.T) {
	var err *ExitError
	if err.Error() != "" {
		t.Errorf("nil ExitError.Error() = %q, want empty", err.Error())
	}
}

func TestExitError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	err := &ExitError{Code: 2, Err: inner}
	if !errors.Is(err, inner) {
		t.Error("Unwrap should return inner error")
	}
}

func TestExitError_Unwrap_Nil(t *testing.T) {
	var err *ExitError
	if err.Unwrap() != nil {
		t.Error("nil ExitError.Unwrap() should return nil")
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"exit error", &ExitError{Code: 3, Err: errors.New("auth")}, 3},
		{"negative code", &ExitError{Code: -1, Err: errors.New("bad")}, 1},
		{"generic error", errors.New("generic"), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExitCode(tt.err)
			if got != tt.want {
				t.Errorf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	got := VersionString()
	if got == "" {
		t.Error("expected non-empty version string")
	}
	// Default is "dev" when not built with ldflags.
	if !strings.Contains(got, "dev") {
		t.Logf("version = %q (may be set by ldflags)", got)
	}
}

func TestWrapParseError_Nil(t *testing.T) {
	if wrapParseError(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestWrapParseError_Generic(t *testing.T) {
	err := errors.New("generic")
	got := wrapParseError(err)
	if got != err {
		t.Errorf("expected same error back, got %v", got)
	}
}
