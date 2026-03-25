package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	noVerify := false
	dryRun := false
	amend := false
	noEdit := false
	allowEmpty := false

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--no-verify", "-nv":
			noVerify = true
		case "--dry-run", "-dr":
			dryRun = true
		case "--amend", "-am":
			amend = true
		case "--no-edit", "-ne":
			noEdit = true
		case "--allow-empty", "-ae":
			allowEmpty = true
		}
	}

	// --no-edit skips the wizard entirely.
	if noEdit {
		if err := StreamGitCommit("", noVerify, amend, noEdit, allowEmpty); err != nil {
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(newModel())

	final, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := final.(model); ok && m.commitMsg != "" {
		if dryRun {
			fmt.Println(m.commitMsg)
			return
		}
		if err := StreamGitCommit(m.commitMsg, noVerify, amend, noEdit, allowEmpty); err != nil {
			os.Exit(1)
		}
	}
}
