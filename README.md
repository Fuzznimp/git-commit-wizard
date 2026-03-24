# git-commit-wizard

Interactive terminal UI for writing [conventional commits](https://www.conventionalcommits.org/).

## Usage

```bash
gcm
```

## Steps

1. **Type** — search and select a commit type. Fuzzy search filters the list with prefix-first priority.
2. **Scope** — optional scope, with suggestions from your git log history. Navigate suggestions with `↑ / ↓`.
3. **Subject** — commit message subject. A character counter tracks the total commit message length against the 72-character limit.

On the last step, pressing `enter` runs the commit immediately with output streamed live to the terminal.

Output format: `type(scope): subject`

## Features

- **Fuzzy search** on commit types — prefix matches are ranked first.
- **Scope suggestions** from your git log, sorted by frequency.
- **Live preview** of the full commit message at every step.
- **Character counter** right-aligned to the 72-character limit, turns red when exceeded.
- **Staged files list** at the bottom — color-coded by status: green `+` for added, orange `~` for modified, red `-` for deleted. Shows a warning when nothing is staged.
- **Commit message displayed** in the UI after the commit runs.

## Keybindings

| Key      | Action                |
| -------- | --------------------- |
| `↑ / ↓`  | Navigate list         |
| `enter`  | Confirm selection     |
| `esc`    | Clear input / go back |
| `ctrl+c` | Quit                  |

## Build

```bash
go build .
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Nerd Fonts](https://www.nerdfonts.com/)
