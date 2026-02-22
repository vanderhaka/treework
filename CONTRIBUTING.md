# Contributing

Thanks for your interest in contributing to treework!

## Getting started

1. Fork and clone the repo
2. Make sure you have Go 1.24+ installed
3. Build: `go build -o treework .`
4. Run: `./treework`

## Making changes

1. Create a branch for your change
2. Make your changes
3. Run `go vet ./...` to check for issues
4. Test manually (no test suite yet — contributions welcome!)
5. Open a pull request

## Guidelines

- Keep it simple. treework is intentionally minimal.
- Follow existing code style and patterns.
- One feature per PR.
- Update the CHANGELOG if your change is user-facing.

## Reporting bugs

Open an issue with:
- What you expected to happen
- What actually happened
- Your OS and Go version

## Feature requests

Open an issue describing the feature and why it would be useful. Keep in mind that treework targets non-technical users who want a simple git worktree workflow.
