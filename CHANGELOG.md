# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.0] - 2025-02-22

### Added

- Interactive menu with arrow key navigation
- Create worktrees with automatic branch creation
- List worktrees with branch and repo info
- Remove worktrees with safety checks for uncommitted changes and unpushed commits
- Remove all worktrees for a repo with per-worktree safety warnings
- Auto-detect and copy `.env` files from main repo
- Package manager detection (npm, yarn, pnpm, bun) with install prompt
- Editor auto-detection (Cursor, VS Code) with `WT_EDITOR` override
- Configurable base folder via settings menu or `DEV_DIR` env var
- Interactive directory browser for setting base folder
- Branch cleanup: auto-delete merged branches, prompt for unmerged
- Keyboard shortcuts: arrow keys, Enter, Esc, Ctrl+C
- Git installation check at startup
