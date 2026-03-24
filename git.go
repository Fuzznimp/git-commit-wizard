package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var conventionalScopeRe = regexp.MustCompile(`^\w+\(([^)]+)\):`)

// GetScopesFromLog runs git log for the current user and extracts scopes
// from conventional commit messages, sorted by frequency descending.
func GetScopesFromLog() []string {
	email, err := gitConfig("user.email")
	if err != nil || email == "" {
		return nil
	}

	out, err := exec.Command("git", "log", "--author="+email, "--format=%s").Output()
	if err != nil {
		return nil
	}

	freq := map[string]int{}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)

		if m := conventionalScopeRe.FindStringSubmatch(line); m != nil {
			freq[m[1]]++
		}
	}

	// Sort by frequency descending.
	type pair struct {
		scope string
		count int
	}

	pairs := make([]pair, 0, len(freq))
	for s, c := range freq {
		pairs = append(pairs, pair{s, c})
	}

	// Simple insertion sort — scope lists are small.
	for i := 1; i < len(pairs); i++ {
		for j := i; j > 0 && pairs[j].count > pairs[j-1].count; j-- {
			pairs[j], pairs[j-1] = pairs[j-1], pairs[j]
		}
	}

	scopes := make([]string, len(pairs))
	for i, p := range pairs {
		scopes[i] = p.scope
	}

	return scopes
}

func gitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", key).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// StagedFile holds a staged file's status and path.
type StagedFile struct {
	Status string
	Path   string
}

// GetStagedFiles returns the list of staged files with their status (A/M/D).
func GetStagedFiles() []StagedFile {
	out, err := exec.Command("git", "diff", "--cached", "--name-status").Output()
	if err != nil {
		return nil
	}

	var files []StagedFile
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			files = append(files, StagedFile{Status: parts[0], Path: parts[1]})
		}
	}

	return files
}

// StreamGitCommit runs `git commit -m msg` with output streamed directly to the terminal.
func StreamGitCommit(msg string) error {
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
