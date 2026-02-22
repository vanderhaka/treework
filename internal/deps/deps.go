package deps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Manager holds the detected package manager info.
type Manager struct {
	Name    string
	Command string
}

// Detect detects the package manager for a project directory.
// Returns nil if no package.json is found.
func Detect(dir string) *Manager {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		return nil
	}

	lockfiles := []struct {
		file string
		name string
	}{
		{"bun.lockb", "bun"},
		{"bun.lock", "bun"},
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
		{"package-lock.json", "npm"},
	}

	for _, lf := range lockfiles {
		if _, err := os.Stat(filepath.Join(dir, lf.file)); err == nil {
			if _, err := exec.LookPath(lf.name); err == nil {
				return &Manager{Name: lf.name, Command: lf.name}
			}
		}
	}

	// Default to npm if package.json exists
	if _, err := exec.LookPath("npm"); err == nil {
		return &Manager{Name: "npm", Command: "npm"}
	}

	return nil
}

// Install runs the package manager install command in the given directory.
func Install(dir string, pm *Manager) error {
	if pm == nil {
		return fmt.Errorf("no package manager provided")
	}
	cmd := exec.Command(pm.Command, "install")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}
