package cmd

import (
	"fmt"
	"os"
	"strings"
)

// Set via ldflags at build time.
var (
	version = "dev"
	commit  = ""
	date    = ""
)

// VersionCmd shows version information.
type VersionCmd struct{}

// Run prints the version string.
func (c *VersionCmd) Run() error {
	_, _ = fmt.Fprintln(os.Stdout, VersionString())
	return nil
}

// VersionString returns a formatted version string.
func VersionString() string {
	var sb strings.Builder
	sb.WriteString(version)

	if commit != "" {
		sb.WriteString(" (")
		sb.WriteString(commit)
		if date != "" {
			sb.WriteString(" ")
			sb.WriteString(date)
		}
		sb.WriteString(")")
	}

	return sb.String()
}
