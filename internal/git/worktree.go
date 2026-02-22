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
		cmd = exec.Command("git", "-C", repoDir, "worktree", "add", wtPath, "-b", branchName)
	} else {
		cmd = exec.Command("git", "-C", repoDir, "worktree", "add", wtPath, branchName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// WorktreeStatus describes the state of a worktree's working directory.
type WorktreeStatus struct {
	HasUncommittedChanges bool // Modified, staged, or untracked files
	HasUnpushedCommits    bool // Commits not pushed to any remote
}

// CheckWorktreeStatus inspects a worktree for unsaved work.
// If git commands fail, assumes the worktree is dirty (fail safe).
func CheckWorktreeStatus(wtPath string) WorktreeStatus {
	var s WorktreeStatus

	// Check for modified, staged, or untracked files
	out, err := exec.Command("git", "-C", wtPath, "status", "--porcelain").Output()
	if err != nil {
		// Git failed — assume dirty to prevent accidental deletion
		s.HasUncommittedChanges = true
		return s
	}
	if len(strings.TrimSpace(string(out))) > 0 {
		s.HasUncommittedChanges = true
	}

	// Check for commits not on any remote
	branch := CurrentBranch(wtPath)
	if branch != "" && branch != "HEAD" {
		out, err = exec.Command("git", "-C", wtPath, "log", branch, "--not", "--remotes", "--oneline").Output()
		if err != nil {
			// Git failed — assume dirty to prevent accidental deletion
			s.HasUnpushedCommits = true
		} else if len(strings.TrimSpace(string(out))) > 0 {
			s.HasUnpushedCommits = true
		}
	}

	return s
}

// IsDirty returns true if the worktree has any unsaved work that would be lost.
func (s WorktreeStatus) IsDirty() bool {
	return s.HasUncommittedChanges || s.HasUnpushedCommits
}

// WorktreeRemove removes a clean worktree. Returns an error if the worktree
// has uncommitted changes (does NOT force).
func WorktreeRemove(repoDir, wtPath string) error {
	return exec.Command("git", "-C", repoDir, "worktree", "remove", wtPath).Run()
}

// WorktreeForceRemove removes a worktree even if it has uncommitted changes.
func WorktreeForceRemove(repoDir, wtPath string) error {
	return exec.Command("git", "-C", repoDir, "worktree", "remove", "--force", wtPath).Run()
}

// WorktreePrune prunes stale worktree references.
func WorktreePrune(repoDir string) {
	exec.Command("git", "-C", repoDir, "worktree", "prune").Run()
}

// WorktreeList returns all worktrees for a repo (excluding the main one).
func WorktreeList(repoDir string) []WorktreeInfo {
	out, err := exec.Command("git", "-C", repoDir, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil
	}

	var worktrees []WorktreeInfo
	var current WorktreeInfo
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	first := true

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "worktree ") {
			if !first && current.Path != "" {
				worktrees = append(worktrees, current)
			}
			first = false
			current = WorktreeInfo{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if strings.HasPrefix(line, "branch ") {
			branch := strings.TrimPrefix(line, "branch ")
			branch = strings.TrimPrefix(branch, "refs/heads/")
			current.Branch = branch
		}
	}

	// Don't forget the last entry, but skip first (main worktree)
	if current.Path != "" && len(worktrees) > 0 {
		worktrees = append(worktrees, current)
	} else if current.Path != "" && !first {
		// This was the only entry (the main worktree) — skip it
	}

	return worktrees
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
			if !strings.Contains(path, "/.git/") {
				dirs = append(dirs, path)
			}
			return fs.SkipDir
		}

		return nil
	})

	return dirs
}
