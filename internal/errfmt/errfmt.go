// Package errfmt formats errors into user-friendly messages with actionable hints.
package errfmt

import (
	"errors"
	"fmt"
	"strings"
)

// Format returns a user-friendly error message.
// It maps known error patterns to helpful hints.
func Format(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()

	// Check wrapped sentinel errors first.
	for _, m := range matchers {
		if m.sentinel != nil && errors.Is(err, m.sentinel) {
			return m.message
		}
	}

	// Fall back to substring matching on the message.
	lower := strings.ToLower(msg)
	for _, m := range matchers {
		if m.sentinel == nil && m.substr != "" && strings.Contains(lower, m.substr) {
			return m.message
		}
	}

	return msg
}

type matcher struct {
	sentinel error
	substr   string
	message  string
}

var matchers = []matcher{
	{
		substr:  "no api key",
		message: "No API key configured. Run: lametric setup",
	},
	{
		substr:  "api key",
		message: "Invalid or missing API key. Run: lametric setup",
	},
	{
		substr:  "unauthorized",
		message: "Authentication failed. Check your API key: lametric setup",
	},
	{
		substr:  "401",
		message: "Authentication failed. Check your API key: lametric setup",
	},
	{
		substr:  "forbidden",
		message: "Access denied. Verify device permissions.",
	},
	{
		substr:  "403",
		message: "Access denied. Verify device permissions.",
	},
	{
		substr:  "not found",
		message: "Resource not found. Check the device name or ID.",
	},
	{
		substr:  "404",
		message: "Resource not found. Check the device name or ID.",
	},
	{
		substr:  "connection refused",
		message: "Cannot reach device. Is it powered on and connected?",
	},
	{
		substr:  "no such host",
		message: "Cannot resolve device address. Check the hostname or IP.",
	},
	{
		substr:  "timeout",
		message: "Request timed out. Is the device reachable on your network?",
	},
	{
		substr:  "tls",
		message: "TLS error connecting to device. Try --insecure if using self-signed certs.",
	},
	{
		substr:  "certificate",
		message: "Certificate error. Try --insecure if using self-signed certs.",
	},
}

// Wrap returns an error with a user-friendly prefix.
func Wrap(err error, action string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", action, err)
}

// Hint appends a hint to the error message.
func Hint(err error, hint string) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%s\n  hint: %s", Format(err), hint)
}
