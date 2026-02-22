package cmd

import (
	"fmt"
	"os"

	"github.com/jamesvanderhaak/wt/internal/config"
	"github.com/jamesvanderhaak/wt/internal/ui"
)

// SetBaseDir runs the shared path-selection flow.
// Returns the chosen directory path or an error.
func SetBaseDir(currentPath string) (string, error) {
	method, err := ui.SelectPathMethod()
	if err != nil {
		return "", err
	}

	var selected string
	switch method {
	case "type":
		selected, err = ui.InputPath(currentPath)
		if err != nil {
			return "", err
		}
	case "browse":
		startDir := currentPath
		if info, serr := os.Stat(startDir); serr != nil || !info.IsDir() {
			startDir, _ = os.UserHomeDir()
		}
		selected, err = ui.BrowseDirectory(startDir)
		if err != nil {
			return "", err
		}
	}

	// Validate the path exists and is a directory
	info, err := os.Stat(selected)
	if err != nil || !info.IsDir() {
		ui.Error(fmt.Sprintf("Not a valid directory: %s", selected))
		return "", fmt.Errorf("invalid directory: %s", selected)
	}

	return selected, nil
}

// doSettings shows the current base folder and lets the user change it.
func doSettings() {
	cfg := config.Load()
	current := config.DevDir()

	fmt.Println()
	ui.Info(fmt.Sprintf("Base folder: %s", current))
	if cfg.BaseDir != "" {
		ui.Muted("(from config file)")
	} else if os.Getenv("DEV_DIR") != "" {
		ui.Muted("(from DEV_DIR env var)")
	} else {
		ui.Muted("(default)")
	}
	fmt.Println()

	selected, err := SetBaseDir(current)
	if err != nil {
		if isAbort(err) {
			return
		}
		return
	}

	cfg.BaseDir = selected
	if err := config.Save(cfg); err != nil {
		ui.Error(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	ui.Success(fmt.Sprintf("Base folder set to %s", selected))
}

// runSettingsInteractive is called from the root menu.
func runSettingsInteractive() {
	doSettings()
}
