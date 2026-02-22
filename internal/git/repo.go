package git

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CurrentRepo returns the git toplevel of the current directory, or empty string if not in a repo.
func CurrentRepo() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// ScanRepos finds all git repos under devDir (maxdepth 5), excluding worktree dirs.
func ScanRepos(devDir string) []string {
	var repos []string
	maxDepth := strings.Count(filepath.Clean(devDir), string(os.PathSeparator)) + 5

	filepath.WalkDir(devDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		depth := strings.Count(filepath.Clean(path), string(os.PathSeparator))
		if depth > maxDepth {
			return fs.SkipDir
		}

		// Skip worktree directories
		if d.IsDir() && strings.Contains(d.Name(), "-worktree-") {
			return fs.SkipDir
		}

		if d.Name() == ".git" && d.IsDir() {
			repos = append(repos, filepath.Dir(path))
			return fs.SkipDir
		}

		return nil
	})

	return repos
}
