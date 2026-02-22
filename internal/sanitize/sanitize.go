package sanitize

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNum = regexp.MustCompile(`[^a-z0-9_-]`)
	multiDash   = regexp.MustCompile(`-{2,}`)
)

// Name sanitizes a worktree/branch name: lowercase, only alphanumeric/hyphens/underscores,
// no leading/trailing hyphens, no consecutive hyphens.
func Name(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = multiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	const maxNameLen = 100
	if len(s) > maxNameLen {
		s = s[:maxNameLen]
		s = strings.TrimRight(s, "-")
	}
	return s
}
