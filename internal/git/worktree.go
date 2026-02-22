package git

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WorktreeInfo holds parsed worktree data.
type WorktreeInfo struct {
	Path   string
	Branch string
}

// WorktreeAdd creates a new worktree. If newBranch is true, creates a new branch.
func WorktreeAdd(repoDir, wtPath, branchName string, newBranch bool) error {
	var cmd *exec.Cmd
	if newBranch {
		cmd = exec.Command("git", "-C", repoDir, "worktree", "add", "-b", branchName, "--", wtPath)
	} else {
		cmd = exec.Command("git", "-C", repoDir, "worktree", "add", "--", wtPath, "--", branchName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// WorktreeRemove removes a worktree without forcing.
func WorktreeRemove(repoDir, wtPath string) error {
	cmd := exec.Command("git", "-C", repoDir, "worktree", "remove", "--", wtPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	msg := strings.TrimSpace(string(out))
	if msg == "" {
		return err
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// WorktreeForceRemove force-removes a worktree.
func WorktreeForceRemove(repoDir, wtPath string) error {
	cmd := exec.Command("git", "-C", repoDir, "worktree", "remove", "--force", "--", wtPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	msg := strings.TrimSpace(string(out))
	if msg == "" {
		return err
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// WorktreePrune prunes stale worktree references.
func WorktreePrune(repoDir string) error {
	return exec.Command("git", "-C", repoDir, "worktree", "prune").Run()
}

// WorktreeList returns all worktrees for a repo (excluding the main one).
func WorktreeList(repoDir string) []WorktreeInfo {
	out, err := exec.Command("git", "-C", repoDir, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil
	}

	var all []WorktreeInfo
	var current WorktreeInfo
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "worktree ") {
			if current.Path != "" {
				all = append(all, current)
			}
			current = WorktreeInfo{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if strings.HasPrefix(line, "branch ") {
			branch := strings.TrimPrefix(line, "branch ")
			branch = strings.TrimPrefix(branch, "refs/heads/")
			current.Branch = branch
		}
	}

	if current.Path != "" {
		all = append(all, current)
	}

	if len(all) > 1 {
		return all[1:]
	}

	return nil
}

// MainWorktreePath returns the path of the main worktree for a repo.
func MainWorktreePath(wtPath string) string {
	out, err := exec.Command("git", "-C", wtPath, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "worktree ") {
			return strings.TrimPrefix(line, "worktree ")
		}
	}
	return ""
}

// CurrentBranch returns the current branch name for a worktree path.
func CurrentBranch(wtPath string) string {
	out, err := exec.Command("git", "-C", wtPath, "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// WorktreePath computes the worktree path: ../{repo}-worktree-{name}
func WorktreePath(repoDir, name string) string {
	repo := filepath.Base(repoDir)
	return filepath.Join(filepath.Dir(repoDir), fmt.Sprintf("%s-worktree-%s", repo, name))
}

// FindWorktreeDirs scans devDir for directories matching *-worktree-* pattern.
func FindWorktreeDirs(devDir string) []string {
	var dirs []string
	maxDepth := strings.Count(filepath.Clean(devDir), string(os.PathSeparator)) + 3

	filepath.WalkDir(devDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		depth := strings.Count(filepath.Clean(path), string(os.PathSeparator))
		if depth > maxDepth {
			return fs.SkipDir
		}

		if d.IsDir() && strings.Contains(d.Name(), "-worktree-") {
			// Skip .git subdirectories
			if !strings.Contains(path, string(os.PathSeparator)+".git"+string(os.PathSeparator)) {
				dirs = append(dirs, path)
			}
			return fs.SkipDir
		}

		return nil
	})

	return dirs
}
