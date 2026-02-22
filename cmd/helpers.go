package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/vanderhaka/treework/internal/config"
	"github.com/vanderhaka/treework/internal/git"
	"github.com/vanderhaka/treework/internal/ui"
)

// resolveRepo returns a repo directory.
// When forceSelect is true (interactive menu), it always shows the project list.
// When false (direct CLI), it tries the current directory first.
func resolveRepo(forceSelect bool) (string, error) {
	if !forceSelect {
		if repo := git.CurrentRepo(); repo != "" {
			return repo, nil
		}
	}

	// Scan base folder for repos
	devDir := requireDevDir()
	if devDir == "" {
		return "", fmt.Errorf("no base folder configured")
	}

	repos := git.ScanRepos(devDir)
	if len(repos) == 0 {
		return "", fmt.Errorf("no git repos found in %s — check your base folder in 'treework settings'", devDir)
	}

	selected, err := ui.SelectRepo(repos)
	if err != nil {
		if isAbort(err) {
			return "", err
		}
		return "", err
	}

	return selected, nil
}

// isAbort checks if an error is a user abort (Escape / Ctrl+C).
func isAbort(err error) bool {
	return errors.Is(err, huh.ErrUserAborted)
}

// handleAbort exits the program cleanly on abort. Use for direct CLI commands only.
func handleAbort(err error) {
	if isAbort(err) {
		fmt.Println()
		ui.Muted("Cancelled.")
		os.Exit(0)
	}
}

// requireDevDir returns the configured base folder, or an empty string with an
// error message if none is set. Call this before any command that needs DevDir.
func requireDevDir() string {
	devDir := config.DevDir()
	if devDir == "" {
		ui.Error("No base folder configured. Run 'treework settings' or set DEV_DIR.")
		return ""
	}
	if _, err := os.Stat(devDir); err != nil {
		ui.Error(fmt.Sprintf("Base folder not found: %s", devDir))
		return ""
	}
	return devDir
}

// resolveWorktreePath resolves a worktree path to an absolute path.
func resolveWorktreePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
