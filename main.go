package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(newModel())

	final, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := final.(model); ok && m.commitMsg != "" {
		if err := StreamGitCommit(m.commitMsg); err != nil {
			os.Exit(1)
		}
	}
}
