package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/jamesvanderhaak/wt/internal/config"
	"github.com/jamesvanderhaak/wt/internal/git"
	"github.com/jamesvanderhaak/wt/internal/ui"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove"},
	Short:   "Remove a worktree",
	Run:     runRm,
}

func runRm(cmd *cobra.Command, args []string) {
	fmt.Println()
	doRm(true)
}

func runRmInteractive(cmd *cobra.Command) {
	doRm(false)
}

func doRm(direct bool) {
	devDir := config.DevDir()
	if _, err := os.Stat(devDir); err != nil {
		ui.Error(fmt.Sprintf("DEV_DIR not found: %s", devDir))
		if direct {
			os.Exit(1)
		}
		return
	}

	dirs := git.FindWorktreeDirs(devDir)
	if len(dirs) == 0 {
		ui.Info("No worktrees found.")
		return
	}

	selected, err := ui.SelectWorktree(dirs)
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

	if selected == ui.BackValue {
		return
	}

	branch := git.CurrentBranch(selected)
	mainDir := git.MainWorktreePath(selected)
	if mainDir == "" {
		ui.Error("Can't find main repo for this worktree.")
		if direct {
			os.Exit(1)
		}
		return
	}
	defaultBranch := git.DefaultBranch(mainDir)

	ui.Info(fmt.Sprintf("Removing: %s (branch: %s)", filepath.Base(selected), branch))

	var removeErr error
	err = spinner.New().
		Title("Removing worktree...").
		Action(func() {
			removeErr = git.WorktreeRemove(mainDir, selected)
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

	if removeErr != nil {
		if strings.Contains(strings.ToLower(removeErr.Error()), "unclean working tree") {
			forceRemove, confirmErr := ui.Confirm(fmt.Sprintf("Force-remove worktree '%s' with uncommitted changes?", filepath.Base(selected)))
			if confirmErr != nil {
				if isAbort(confirmErr) {
					if direct {
						handleAbort(confirmErr)
					}
					return
				}
				ui.Error(confirmErr.Error())
				if direct {
					os.Exit(1)
				}
				return
			}
			if forceRemove {
				removeErr = git.WorktreeForceRemove(mainDir, selected)
			}
		}
	}

	if removeErr != nil {
		_ = git.WorktreePrune(mainDir)
		ui.Error("Failed to remove worktree.")
		if direct {
			os.Exit(1)
		}
		return
	}

	_ = git.WorktreePrune(mainDir)
	ui.Success("Removed worktree")

	if branch != "" && branch != "HEAD" && branch != defaultBranch {
		if git.IsBranchMerged(mainDir, branch) {
			if err := git.DeleteBranch(mainDir, branch); err == nil {
				ui.Success(fmt.Sprintf("Deleted merged branch '%s'", branch))
			} else {
				ui.Warn(fmt.Sprintf("Could not delete branch '%s': %v", branch, err))
			}
		} else {
			ui.Warn(fmt.Sprintf("Branch '%s' is not merged", branch))
			forceDelete, err := ui.ConfirmForceDelete(branch)
			if err != nil {
				if isAbort(err) {
					if direct {
						handleAbort(err)
					}
					return
				} else {
					ui.Warn(fmt.Sprintf("Prompt error: %v — keeping branch", err))
					return
				}
			}
			if forceDelete {
				if err := git.ForceDeleteBranch(mainDir, branch); err == nil {
					ui.Success(fmt.Sprintf("Force deleted branch '%s'", branch))
				} else {
					ui.Error(fmt.Sprintf("Failed to delete branch '%s'", branch))
				}
			} else {
				ui.Muted(fmt.Sprintf("Kept branch '%s'", branch))
			}
		}
	}
}
