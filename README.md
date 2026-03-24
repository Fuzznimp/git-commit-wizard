# git-commit-wizard

Interactive terminal UI for writing [conventional commits](https://www.conventionalcommits.org/).

Three steps:

1. **Type** — search and select a commit type (`feat`, `fix`, `chore`, etc.).
2. **Scope** — optional scope, with suggestions from your git history.
3. **Subject** — commit message subject.

Press `enter` to confirm each step. The commit runs immediately on the last step, with git output streamed live to the terminal.

Output format: `type(scope): subject`

## Keybindings

| Key     | Action                |
| ------- | --------------------- |
| `↑ / ↓` | Navigate list         |
| `enter` | Confirm selection     |
| `esc`   | Clear input / go back |

## Build

```bash
go build .
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Nerd Fonts](https://www.nerdfonts.com/)
