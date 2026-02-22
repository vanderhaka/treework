package editor

import (
	"fmt"
	"os/exec"

	"github.com/jamesvanderhaak/wt/internal/config"
)

// Open opens the given path in the preferred editor.
// Priority: WT_EDITOR env var > cursor > code > open (macOS).
func Open(path string) error {
	if ed := config.Editor(); ed != "" {
		return run(ed, path)
	}

	editors := []struct {
		cmd  string
		args []string
	}{
		{"cursor", []string{"--new-window", path}},
		{"code", []string{"-n", path}},
		{"open", []string{path}},
	}

	for _, e := range editors {
		if _, err := exec.LookPath(e.cmd); err == nil {
			return exec.Command(e.cmd, e.args...).Start()
		}
	}

	return fmt.Errorf("no editor found: install cursor, code, or set WT_EDITOR")
}

func run(editor, path string) error {
	switch editor {
	case "cursor":
		return exec.Command(editor, "--new-window", path).Start()
	case "code":
		return exec.Command(editor, "-n", path).Start()
	default:
		return exec.Command(editor, path).Start()
	}
}
