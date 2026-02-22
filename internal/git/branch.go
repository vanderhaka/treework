package git

import (
	"os/exec"
	"strings"
)

// BranchExists checks if a local branch exists in the given repo.
func BranchExists(repoDir, branch string) bool {
	err := exec.Command("git", "-C", repoDir, "show-ref", "--verify", "--quiet", "--", "refs/heads/"+branch).Run()
	return err == nil
}

// DefaultBranch returns the default branch for the repo (main or master).
func DefaultBranch(repoDir string) string {
	// Try origin/HEAD first
	out, err := exec.Command("git", "-C", repoDir, "symbolic-ref", "--quiet", "--short", "refs/remotes/origin/HEAD").Output()
	if err == nil {
		d := strings.TrimSpace(string(out))
		d = strings.TrimPrefix(d, "origin/")
		if d != "" {
			return d
		}
	}

	// Check for main
	if err := exec.Command("git", "-C", repoDir, "show-ref", "--verify", "--quiet", "refs/heads/main").Run(); err == nil {
		return "main"
	}

	// Check for master
	if err := exec.Command("git", "-C", repoDir, "show-ref", "--verify", "--quiet", "refs/heads/master").Run(); err == nil {
		return "master"
	}

	return "main"
}

// IsBranchMerged checks if branch is merged into the default branch.
func IsBranchMerged(repoDir, branch string) bool {
	base := DefaultBranch(repoDir)
	err := exec.Command("git", "-C", repoDir, "merge-base", "--is-ancestor", "--", branch, base).Run()
	return err == nil
}

// DeleteBranch deletes a local branch (soft delete with -d).
func DeleteBranch(repoDir, branch string) error {
	return exec.Command("git", "-C", repoDir, "branch", "-d", "--", branch).Run()
}

// ForceDeleteBranch deletes a local branch with -D (force).
func ForceDeleteBranch(repoDir, branch string) error {
	return exec.Command("git", "-C", repoDir, "branch", "-D", "--", branch).Run()
}
