package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// BackValue is the sentinel value returned when the user picks "← Back".
const BackValue = "__back__"

// keymap returns a custom huh keymap with Escape and left arrow mapped to quit (back).
func keymap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc", "left"))
	km.Select.Submit = key.NewBinding(key.WithKeys("enter", "right"), key.WithHelp("enter/→", "select"))
	return km
}

// defaultKeymap returns the default keymap (only ctrl+c quits).
// Used for the main menu where Escape should exit, not go back.
func defaultKeymap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
	km.Select.Submit = key.NewBinding(key.WithKeys("enter", "right"), key.WithHelp("enter/→", "select"))
	return km
}

// confirmKeymap returns a keymap for yes/no confirms.
// Left/right are NOT overridden so they toggle between Yes and No naturally.
func confirmKeymap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
	return km
}

// runField wraps a single huh field in a form with the back-enabled keymap.
func runField(field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).WithKeyMap(keymap()).Run()
}

// runFieldConfirm wraps a confirm field with a keymap safe for yes/no toggles.
func runFieldConfirm(field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).WithKeyMap(confirmKeymap()).Run()
}

// runFieldDefault wraps a single huh field in a form with the default keymap.
func runFieldDefault(field huh.Field) error {
	return huh.NewForm(huh.NewGroup(field)).WithKeyMap(defaultKeymap()).Run()
}

// SelectRepo prompts the user to pick a repo from a list.
func SelectRepo(repos []string) (string, error) {
	var opts []huh.Option[string]
	for _, r := range repos {
		label := filepath.Base(r)
		opts = append(opts, huh.NewOption(label, r))
	}

	var selected string
	field := huh.NewSelect[string]().
		Title("Select a project").
		Options(opts...).
		Value(&selected)

	err := runField(field)
	return selected, err
}

// InputName prompts the user to enter a worktree name.
func InputName() (string, error) {
	var name string
	field := huh.NewInput().
		Title("Worktree name").
		Placeholder("feature-name").
		Value(&name)

	err := runField(field)
	return name, err
}

// WorktreeDisplay holds display info for a worktree in the selector.
type WorktreeDisplay struct {
	Path   string
	Branch string
	Repo   string
}

// SelectWorktree prompts the user to pick a worktree from a list.
// Returns BackValue if the user picks "← Back".
func SelectWorktree(dirs []string) (string, error) {
	opts := []huh.Option[string]{
		huh.NewOption(MutedStyle.Render("← Back"), BackValue),
	}
	for _, d := range dirs {
		label := filepath.Base(d)
		opts = append(opts, huh.NewOption(label, d))
	}

	var selected string
	field := huh.NewSelect[string]().
		Title("Select a worktree").
		Options(opts...).
		Value(&selected)

	err := runField(field)
	return selected, err
}

// SelectWorktreeDetailed prompts the user to pick a worktree, showing branch and repo info.
// Returns BackValue if the user picks "← Back".
func SelectWorktreeDetailed(items []WorktreeDisplay) (string, error) {
	opts := []huh.Option[string]{
		huh.NewOption(MutedStyle.Render("← Back"), BackValue),
	}
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
	field := huh.NewSelect[string]().
		Title("Worktrees").
		Options(opts...).
		Value(&selected)

	err := runField(field)
	return selected, err
}

// ConfirmOpen prompts whether to open the selected worktree in the editor.
func ConfirmOpen(name string) (bool, error) {
	return Confirm(fmt.Sprintf("Open %s in editor?", name))
}

// SelectAction prompts the user to pick an action from the interactive menu.
func SelectAction() (string, error) {
	var action string
	field := huh.NewSelect[string]().
		Title("What would you like to do?").
		Options(
			huh.NewOption("Create new worktree", "new"),
			huh.NewOption("List worktrees", "ls"),
			huh.NewOption("Remove a worktree", "rm"),
			huh.NewOption("Remove ALL worktrees for a repo", "clear"),
			huh.NewOption(MutedStyle.Render("Settings"), "settings"),
			huh.NewOption(MutedStyle.Render("Quit"), "quit"),
		).
		Value(&action)

	err := runFieldDefault(field)
	return action, err
}

// Confirm prompts the user for a yes/no confirmation.
// Uses a dedicated keymap so left/right toggle Yes/No instead of acting as back/submit.
func Confirm(title string) (bool, error) {
	var confirmed bool
	field := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed)

	err := runFieldConfirm(field)
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

// SelectPathMethod prompts the user to choose how to set the base folder path.
func SelectPathMethod() (string, error) {
	var method string
	field := huh.NewSelect[string]().
		Title("How would you like to set the base folder?").
		Options(
			huh.NewOption("Type a path", "type"),
			huh.NewOption("Browse directories", "browse"),
		).
		Value(&method)

	err := runField(field)
	return method, err
}

// InputPath prompts the user to type a directory path.
func InputPath(defaultPath string) (string, error) {
	value := defaultPath
	field := huh.NewInput().
		Title("Base folder path").
		Value(&value)

	err := runField(field)
	return value, err
}

// BrowseDirectory presents an interactive directory browser starting at startDir.
// Returns the chosen directory path or an error.
func BrowseDirectory(startDir string) (string, error) {
	current := startDir
	for {
		entries, err := os.ReadDir(current)
		if err != nil {
			return "", fmt.Errorf("cannot read directory: %w", err)
		}

		opts := []huh.Option[string]{
			huh.NewOption(SuccessStyle.Render("✓ Use this folder"), "__select__"),
		}

		// Add parent directory option unless at filesystem root
		if current != "/" {
			opts = append(opts, huh.NewOption(MutedStyle.Render(".."), "__parent__"))
		}

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			opts = append(opts, huh.NewOption(name, filepath.Join(current, name)))
		}

		var selected string
		field := huh.NewSelect[string]().
			Title(current).
			Options(opts...).
			Value(&selected)

		if err := runField(field); err != nil {
			return "", err
		}

		switch selected {
		case "__select__":
			return current, nil
		case "__parent__":
			current = filepath.Dir(current)
		default:
			current = selected
		}
	}
}
