package ui

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

// SelectRepo prompts the user to pick a repo from a list.
func SelectRepo(repos []string) (string, error) {
	var opts []huh.Option[string]
	for _, r := range repos {
		label := filepath.Base(r)
		opts = append(opts, huh.NewOption(label, r))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Select a project").
		Options(opts...).
		Value(&selected).
		Run()

	return selected, err
}

// InputName prompts the user to enter a worktree name.
func InputName() (string, error) {
	var name string
	err := huh.NewInput().
		Title("Worktree name").
		Placeholder("feature-name").
		Value(&name).
		Run()

	return name, err
}

// WorktreeDisplay holds display info for a worktree in the selector.
type WorktreeDisplay struct {
	Path   string
	Branch string
	Repo   string
}

// SelectWorktree prompts the user to pick a worktree from a list.
func SelectWorktree(dirs []string) (string, error) {
	var opts []huh.Option[string]
	for _, d := range dirs {
		label := filepath.Base(d)
		opts = append(opts, huh.NewOption(label, d))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Select a worktree").
		Options(opts...).
		Value(&selected).
		Run()

	return selected, err
}

// SelectWorktreeDetailed prompts the user to pick a worktree, showing branch and repo info.
func SelectWorktreeDetailed(items []WorktreeDisplay) (string, error) {
	var opts []huh.Option[string]
	for _, item := range items {
		label := filepath.Base(item.Path)
		if item.Branch != "" {
			label += MutedStyle.Render("  (" + item.Branch + ")")
		}
		if item.Repo != "" {
			label += MutedStyle.Render("  " + item.Repo)
		}
		opts = append(opts, huh.NewOption(label, item.Path))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Worktrees").
		Options(opts...).
		Value(&selected).
		Run()

	return selected, err
}

// ConfirmOpen prompts whether to open the selected worktree in the editor.
func ConfirmOpen(name string) (bool, error) {
	return Confirm(fmt.Sprintf("Open %s in editor?", name))
}

// SelectAction prompts the user to pick an action from the interactive menu.
func SelectAction() (string, error) {
	var action string
	err := huh.NewSelect[string]().
		Title("What would you like to do?").
		Options(
			huh.NewOption("Create new worktree", "new"),
			huh.NewOption("List worktrees", "ls"),
			huh.NewOption("Remove a worktree", "rm"),
			huh.NewOption("Remove ALL worktrees for a repo", "clear"),
		).
		Value(&action).
		Run()

	return action, err
}

// Confirm prompts the user for a yes/no confirmation.
func Confirm(title string) (bool, error) {
	var confirmed bool
	err := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()

	return confirmed, err
}

// ConfirmInstall prompts whether to install dependencies.
func ConfirmInstall(pmName string) (bool, error) {
	return Confirm(fmt.Sprintf("Install dependencies with %s?", pmName))
}

// ConfirmForceDelete prompts whether to force-delete an unmerged branch.
func ConfirmForceDelete(branch string) (bool, error) {
	return Confirm(fmt.Sprintf("Force delete unmerged branch '%s'?", branch))
}
