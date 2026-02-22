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
	doClear(true)
}

func doClear(direct bool) {
	fmt.Println()

	// 1. Resolve repo — interactive menu always shows the project list
	repoDir, err := resolveRepo(!direct)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
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

	// 4. Confirm
	confirmed, err := ui.Confirm(fmt.Sprintf("Remove all %d worktrees?", len(worktrees)))
	if err != nil {
		handleAbort(err)
		ui.Error(err.Error())
		os.Exit(1)
	}
	if !confirmed {
		ui.Muted("Cancelled.")
		fmt.Println()
		return
	}

	// 5. Remove each worktree
	err = spinner.New().
		Title("Removing worktrees...").
		Action(func() {
			for _, wt := range worktrees {
				git.WorktreeRemove(repoDir, wt.Path)

				// Auto-delete merged branches
				if wt.Branch != "" && wt.Branch != "main" && wt.Branch != "master" {
					if git.IsBranchMerged(repoDir, wt.Branch) {
						git.DeleteBranch(repoDir, wt.Branch)
					}
				}
			}
			git.WorktreePrune(repoDir)
		}).
		Run()

	if err != nil {
		handleAbort(err)
		ui.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println()
	ui.Success("All worktrees cleared")
	ui.Muted("Merged branches were auto-deleted")
	fmt.Println()
}
