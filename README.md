# treework

A friendly CLI for managing git worktrees. Create, list, and remove worktrees without memorising git commands.

## Why

Git worktrees let you work on multiple branches at the same time, each in its own folder. But the built-in commands are verbose and easy to get wrong. treework wraps them in an interactive menu so you can stay focused on your code.

## Features

- **Interactive menu** — arrow keys and Enter, no commands to remember
- **Auto-detects your projects** — scans your dev folder for git repos
- **Installs dependencies** — detects npm/yarn/pnpm/bun and offers to install after creation
- **Copies `.env` files** — carries over environment config from the main repo
- **Opens your editor** — launches Cursor, VS Code, or your preferred editor
- **Safety checks on removal** — warns you before deleting worktrees with uncommitted changes or unpushed commits
- **Branch cleanup** — auto-deletes merged branches, asks before force-deleting unmerged ones

## Install

### From source

```sh
go install github.com/vanderhaka/treework@latest
```

### Build locally

```sh
git clone https://github.com/vanderhaka/treework.git
cd treework
go build -o treework .
```

Move the binary somewhere on your `$PATH`:

```sh
mv treework /usr/local/bin/
```

## Quick start

```sh
treework
```

That's it. The interactive menu walks you through everything:

```
treework — git worktree manager

┃ What would you like to do?
┃ > Create new worktree
┃   List worktrees
┃   Remove a worktree
┃   Remove ALL worktrees for a repo
┃   Settings
┃   Quit
```

### Direct commands

```sh
treework new feature-auth    # Create a worktree
treework ls                  # List and open worktrees
treework rm                  # Remove a worktree (with safety checks)
treework clear               # Remove all worktrees for a repo
treework settings            # Change your base folder
treework version             # Print version
```

## Configuration

### Base folder

On first run, treework asks where your git repos live. This is saved to `~/.config/treework/config.json`.

You can also set it via environment variable:

```sh
export DEV_DIR=~/projects
```

Priority: `DEV_DIR` env var > config file (no default — you must set one)

### Editor

treework auto-detects Cursor and VS Code. To override:

```sh
export WT_EDITOR=code
```

## How it works

When you create a worktree called `feature-auth` in a repo called `my-app`:

1. Creates `my-app-worktree-feature-auth/` next to your repo
2. Checks out a new branch called `feature-auth`
3. Copies any `.env` files from the main repo
4. Offers to install dependencies
5. Opens the folder in your editor

When you remove a worktree:

1. Checks for uncommitted changes and unpushed commits
2. If unsaved work is found, shows a warning and asks for confirmation
3. Removes the worktree folder
4. Auto-deletes the branch if it's been merged
5. Asks before force-deleting unmerged branches

## Keyboard shortcuts

| Key | Action |
|---|---|
| `↑` `↓` | Navigate |
| `Enter` or `→` | Select |
| `Esc` or `←` | Go back |
| `/` | Filter list |
| `Ctrl+C` | Quit |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
