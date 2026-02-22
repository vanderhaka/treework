package config

import (
	"os"
	"path/filepath"
)

// DevDir returns the root development directory.
// Uses DEV_DIR env var, falling back to ~/Desktop/Development.
func DevDir() string {
	if d := os.Getenv("DEV_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Desktop", "Development")
}

// Editor returns the preferred editor command.
// Uses WT_EDITOR env var, with no default (caller handles detection).
func Editor() string {
	return os.Getenv("WT_EDITOR")
}
