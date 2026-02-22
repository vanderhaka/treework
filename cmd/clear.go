package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh/spinner"
	"github.com/jamesvanderhaak/wt/internal/git"
	"github.com/jamesvanderhaak/wt/internal/ui"
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
	fmt.Println()
	doClear(true)
}

func doClear(direct bool) {

	// 1. Resolve repo
	repoDir, err := resolveRepo()
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
	defaultBranch := git.DefaultBranch(repoDir)

	// 3. Display list
	ui.Info(fmt.Sprintf("Worktrees for %s:", ui.BoldStyle.Render(repoName)))
	for _, wt := range worktrees {
		ui.Muted(fmt.Sprintf("%s  %s", filepath.Base(wt.Path), ui.MutedStyle.Render("("+wt.Branch+")")))
	}
	fmt.Println()

	// 4. Confirm
	confirmed, err := ui.Confirm(fmt.Sprintf("Remove all %d worktrees?", len(worktrees)))
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

	// 5. Remove each worktree
	errs := []error{}
	err = spinner.New().
		Title("Removing worktrees...").
		Action(func() {
			for _, wt := range worktrees {
				if err := git.WorktreeRemove(repoDir, wt.Path); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", filepath.Base(wt.Path), err))
					continue
				}

				// Auto-delete merged branches
				if wt.Branch != "" && wt.Branch != defaultBranch {
					if git.IsBranchMerged(repoDir, wt.Branch) {
						git.DeleteBranch(repoDir, wt.Branch)
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
	if len(errs) > 0 {
		fmt.Println()
		ui.Error("Some worktrees could not be removed:")
		for _, removeErr := range errs {
			ui.Error(fmt.Sprintf("- %s", removeErr))
		}
		return
	}

	fmt.Println()
	ui.Success("All worktrees cleared")
	ui.Muted("Merged branches were auto-deleted")
	fmt.Println()
}
