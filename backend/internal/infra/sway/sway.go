package sway

import (
	"context"
	"fmt"
	"os/exec"
)

// CreateOutput tells Sway to create a new headless output and sets its resolution.
func CreateOutput(ctx context.Context, outputName, resolution string) error {
	// Create output
	cmd := exec.CommandContext(ctx, "swaymsg", "create_output", outputName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create output %s: %w", outputName, err)
	}

	// Set resolution and enable it
	modeCmd := exec.CommandContext(ctx, "swaymsg", "output", outputName, "resolution", resolution, "position", "0,0")
	if err := modeCmd.Run(); err != nil {
		return fmt.Errorf("failed to set resolution for %s: %w", outputName, err)
	}

	return nil
}

// DestroyOutput unplugs the headless output from Sway.
func DestroyOutput(ctx context.Context, outputName string) error {
	cmd := exec.CommandContext(ctx, "swaymsg", "output", outputName, "unplug")
	if err := cmd.Run(); err != nil {
		// It might be 'disable' depending on sway version, let's try unplug first
		return fmt.Errorf("failed to destroy output %s: %w", outputName, err)
	}
	return nil
}
