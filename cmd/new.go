package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh/spinner"
	"github.com/vanderhaka/treework/internal/deps"
	"github.com/vanderhaka/treework/internal/editor"
	"github.com/vanderhaka/treework/internal/env"
	"github.com/vanderhaka/treework/internal/git"
	"github.com/vanderhaka/treework/internal/sanitize"
	"github.com/vanderhaka/treework/internal/ui"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new worktree",
	Args:  cobra.MaximumNArgs(1),
	Run:   runNew,
}

func runNewInteractive(cmd *cobra.Command) {
	doNew(nil, false)
}

func runNew(cmd *cobra.Command, args []string) {
	fmt.Println()
	doNew(args, true)
}

func doNew(args []string, direct bool) {
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

	// 2. Get name (with retry for invalid names and existing worktrees)
	var name, resolved string
	repoName := filepath.Base(repoDir)

	if len(args) > 0 {
		name = sanitize.Name(args[0])
		if name == "" {
			ui.Error("Invalid name — use letters, numbers, hyphens, or underscores.")
			if direct {
				os.Exit(1)
			}
			return
		}
		resolved = resolveWorktreePath(git.WorktreePath(repoDir, name))
		if _, err := os.Stat(resolved); err == nil {
			ui.Info(fmt.Sprintf("'%s' already exists — opening it instead.", name))
			if err := editor.Open(resolved); err != nil {
				ui.Warn(fmt.Sprintf("Could not open editor: %v", err))
			}
			return
		}
	} else {
		for {
			name, err = ui.InputName()
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

			name = sanitize.Name(name)
			if name == "" {
				ui.Warn("Invalid name — use letters, numbers, hyphens, or underscores. Try again.")
				continue
			}

			resolved = resolveWorktreePath(git.WorktreePath(repoDir, name))
			if _, err := os.Stat(resolved); err == nil {
				ui.Warn(fmt.Sprintf("'%s' already exists. Pick a different name.", name))
				continue
			}
			break
		}
	}

	// 3. Create worktree (with spinner)
	branchExists := git.BranchExists(repoDir, name)
	var addErr error

	err = spinner.New().
		Title(fmt.Sprintf("Creating %s/%s...", repoName, name)).
		Action(func() {
			addErr = git.WorktreeAdd(repoDir, resolved, name, !branchExists)
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
	if addErr != nil {
		ui.Error(fmt.Sprintf("Failed to create worktree: %v", addErr))
		if direct {
			os.Exit(1)
		}
		return
	}

	// 7. Copy .env files
	copied, _ := env.CopyEnvFiles(repoDir, resolved)
	if len(copied) > 0 {
		ui.Muted(fmt.Sprintf("Copied %d env file(s)", len(copied)))
	}

	// 8. Detect package manager → prompt to install deps
	if pm := deps.Detect(resolved); pm != nil {
		install, err := ui.ConfirmInstall(pm.Name)
		if err != nil {
			if isAbort(err) {
				if direct {
					handleAbort(err)
				}
				ui.Muted("Skipped dependency install")
			}
		} else if install {
			var installErr error
			err = spinner.New().
				Title(fmt.Sprintf("Installing dependencies with %s...", pm.Name)).
				Action(func() {
					installErr = deps.Install(resolved, pm)
				}).
				Context(context.Background()).
				Run()

			if err != nil && !isAbort(err) {
				ui.Warn("Install failed. You can run it later inside the folder.")
			} else if installErr != nil {
				ui.Warn("Install failed. You can run it later inside the folder.")
			} else if err == nil {
				ui.Success("Dependencies installed")
			}
		}
	}

	// 9. Open in editor
	if err := editor.Open(resolved); err != nil {
		ui.Warn(fmt.Sprintf("Could not open editor: %v", err))
	}

	// Print success
	fmt.Println()
	ui.Success(fmt.Sprintf("Ready: %s/%s", repoName, name))
	ui.Muted(resolved)
}
