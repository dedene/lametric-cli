package cmd

import (
	"context"
	"fmt"
)

// DismissCmd dismisses a notification by ID (or the current one).
type DismissCmd struct {
	ID string `arg:"" optional:"" help:"Notification ID (omit to dismiss current)"`
}

// Run executes the dismiss command.
func (c *DismissCmd) Run(flags *RootFlags) error {
	client, _, err := resolveClient(flags)
	if err != nil {
		return err
	}

	path := "/api/v2/device/notifications/current"
	if c.ID != "" {
		path = "/api/v2/device/notifications/" + c.ID
	}

	if err := client.Delete(context.Background(), path); err != nil {
		return fmt.Errorf("dismiss notification: %w", err)
	}

	f := newFormatter(flags)

	type dismissResult struct {
		Success bool   `json:"success"`
		ID      string `json:"id,omitempty"`
	}

	id := c.ID
	if id == "" {
		id = "current"
	}

	return f.OutputSingle(dismissResult{Success: true, ID: id}, [][2]string{
		{"Status", "dismissed"},
		{"ID", id},
	})
}
