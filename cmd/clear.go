package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh/spinner"
	"github.com/vanderhaka/treework/internal/git"
	"github.com/vanderhaka/treework/internal/ui"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove ALL worktrees for a repo",
	Run:   runClear,
}

// runClearInteractive is called from the root menu loop.
func runClearInteractive(cmd *cobra.Command) {
	doClear(false)
}

func runClear(cmd *cobra.Command, args []string) {
	doClear(true)
}

func doClear(direct bool) {
	fmt.Println()

	// 1. Resolve repo — interactive menu always shows the project list
	repoDir, err := resolveRepo(!direct)
	if err != nil {
		if isAbort(err) {
			if direct {
				handleAbort(err)
			}
			return
		}
		ui.Error(err.Error())
		if direct {
			os.Exit(1)
		}
		return
	}

	repoName := filepath.Base(repoDir)

	// 2. List worktrees (excluding main)
	worktrees := git.WorktreeList(repoDir)
	if len(worktrees) == 0 {
		ui.Info("No worktrees to remove.")
		fmt.Println()
		return
	}

	// 3. Display list
	ui.Info(fmt.Sprintf("Worktrees for %s:", ui.BoldStyle.Render(repoName)))
	for _, wt := range worktrees {
		ui.Muted(fmt.Sprintf("%s  %s", filepath.Base(wt.Path), ui.MutedStyle.Render("("+wt.Branch+")")))
	}
	fmt.Println()

	// 4. Safety check: identify dirty worktrees
	type worktreeCheck struct {
		info   git.WorktreeInfo
		status git.WorktreeStatus
	}
	var dirty []worktreeCheck
	for _, wt := range worktrees {
		s := git.CheckWorktreeStatus(wt.Path)
		entry := worktreeCheck{info: wt, status: s}
		if s.IsDirty() {
			dirty = append(dirty, entry)
		}
	}

	// 5. Show dirty worktree warnings
	if len(dirty) > 0 {
		fmt.Println()
		ui.Warn(fmt.Sprintf("%d worktree(s) have unsaved work:", len(dirty)))
		for _, d := range dirty {
			reasons := ""
			if d.status.HasUncommittedChanges && d.status.HasUnpushedCommits {
				reasons = "uncommitted changes + unpushed commits"
			} else if d.status.HasUncommittedChanges {
				reasons = "uncommitted changes"
			} else {
				reasons = "unpushed commits"
			}
			ui.Muted(fmt.Sprintf("  • %s (%s) — %s", filepath.Base(d.info.Path), d.info.Branch, reasons))
		}
		fmt.Println()
	}

	// 6. Confirm removal
	var confirmed bool
	if len(dirty) > 0 {
		confirmed, err = ui.Confirm(fmt.Sprintf("Remove all %d worktrees? Unsaved work will be permanently lost", len(worktrees)))
	} else {
		confirmed, err = ui.Confirm(fmt.Sprintf("Remove all %d worktrees?", len(worktrees)))
	}
	if err != nil {
		if isAbort(err) {
			if direct {
				handleAbort(err)
			}
			return
		}
		ui.Error(err.Error())
		if direct {
			os.Exit(1)
		}
		return
	}
	if !confirmed {
		ui.Muted("Cancelled.")
		fmt.Println()
		return
	}

	// 7. Remove each worktree, tracking failures and unmerged branches
	var failed []string
	var unmergedBranches []string

	err = spinner.New().
		Title("Removing worktrees...").
		Action(func() {
			for _, wt := range worktrees {
				// Use force only for dirty worktrees (user already confirmed)
				status := git.CheckWorktreeStatus(wt.Path)
				var removeErr error
				if status.IsDirty() {
					removeErr = git.WorktreeForceRemove(repoDir, wt.Path)
				} else {
					removeErr = git.WorktreeRemove(repoDir, wt.Path)
				}

				if removeErr != nil {
					failed = append(failed, filepath.Base(wt.Path))
					continue
				}

				// Branch cleanup
				if wt.Branch != "" && wt.Branch != "main" && wt.Branch != "master" {
					if git.IsBranchMerged(repoDir, wt.Branch) {
						git.DeleteBranch(repoDir, wt.Branch)
					} else {
						unmergedBranches = append(unmergedBranches, wt.Branch)
					}
				}
			}
			git.WorktreePrune(repoDir)
		}).
		Run()

	if err != nil {
		if isAbort(err) {
			if direct {
				handleAbort(err)
			}
			return
		}
		ui.Error(err.Error())
		if direct {
			os.Exit(1)
		}
		return
	}

	fmt.Println()

	// Report failures
	if len(failed) > 0 {
		ui.Warn(fmt.Sprintf("Failed to remove %d worktree(s):", len(failed)))
		for _, f := range failed {
			ui.Muted(fmt.Sprintf("  • %s", f))
		}
	}

	removed := len(worktrees) - len(failed)
	if removed > 0 {
		ui.Success(fmt.Sprintf("Removed %d worktree(s)", removed))
		ui.Muted("Merged branches were auto-deleted")
	}

	// Handle unmerged branches
	if len(unmergedBranches) > 0 {
		fmt.Println()
		ui.Warn(fmt.Sprintf("%d branch(es) are not merged:", len(unmergedBranches)))
		for _, b := range unmergedBranches {
			ui.Muted(fmt.Sprintf("  • %s", b))
		}
		fmt.Println()
		forceDelete, confirmErr := ui.Confirm("Force delete all unmerged branches?")
		if confirmErr != nil {
			if isAbort(confirmErr) {
				if direct {
					handleAbort(confirmErr)
				}
				ui.Muted("Kept unmerged branches")
				return
			}
		}
		if forceDelete {
			for _, b := range unmergedBranches {
				if err := git.ForceDeleteBranch(repoDir, b); err == nil {
					ui.Success(fmt.Sprintf("Deleted branch '%s'", b))
				} else {
					ui.Warn(fmt.Sprintf("Failed to delete branch '%s'", b))
				}
			}
		} else {
			ui.Muted("Kept unmerged branches")
		}
	}

	fmt.Println()
}
