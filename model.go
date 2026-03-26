package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Steps
const (
	stepType    = iota
	stepScope   = iota
	stepSubject = iota
	stepDone    = iota
)

const commitLimit = 72

var commitTypes = []string{
	"feat", "fix", "chore", "docs", "refactor",
	"test", "style", "perf", "ci", "build", "revert",
}

var title = " git-commit-wizard "

// Colors
var colors = struct {
	orange lipgloss.Color
	green  lipgloss.Color
	red    lipgloss.Color
	dim    lipgloss.Color
}{
	orange: "#e78a4e",
	green:  "#a9b665",
	red:    "#ea6962",
	dim:    "#665c54",
}

// Styles
var (
	selectedStyle = lipgloss.NewStyle().
			Foreground(colors.orange).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(colors.dim)

	labelStyle = lipgloss.NewStyle().
			Foreground(colors.orange).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(colors.green).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(colors.green)

	errorStyle = lipgloss.NewStyle().
			Foreground(colors.red)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.dim).
			Padding(0, 1)
)

type model struct {
	step int

	// stepType
	typeInput     textinput.Model
	filteredTypes []string
	typeIdx       int

	// stepScope
	scopeInput     textinput.Model
	allScopes      []string
	filteredScopes []string
	scopeIdx       int

	// stepSubject
	subjectInput textinput.Model

	// results
	commitType    string
	commitScope   string
	commitSubject string
	commitMsg     string

	stagedFiles []StagedFile
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "type"
	ti.Prompt = "󰍉 "
	ti.CharLimit = 20
	ti.Focus()

	si := textinput.New()
	si.Placeholder = "scope - optional"
	si.Prompt = "󰍉 "
	si.CharLimit = 64

	sub := textinput.New()
	sub.Placeholder = "description"
	sub.Prompt = "󰍉 "
	sub.CharLimit = 100

	scopes := GetScopesFromLog()

	return model{
		step:           stepType,
		typeInput:      ti,
		filteredTypes:  commitTypes,
		typeIdx:        0,
		scopeInput:     si,
		subjectInput:   sub,
		allScopes:      scopes,
		filteredScopes: scopes,
		stagedFiles:    GetStagedFiles(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.step {
		case stepType:
			return m.updateType(msg)
		case stepScope:
			return m.updateScope(msg)
		case stepSubject:
			return m.updateSubject(msg)
		}
	}

	return m, nil
}

// updateType handles key events on the type picker step.
func (m model) updateType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc", "ctrl+h":
		if m.typeInput.Value() == "" {
			return m, tea.Quit
		}
		m.typeInput.SetValue("")
		m.filteredTypes = commitTypes
		m.typeIdx = 0
		return m, nil

	case "down", "ctrl+j":
		if m.typeIdx < len(m.filteredTypes)-1 {
			m.typeIdx++
		}
		return m, nil

	case "up", "ctrl+k":
		if m.typeIdx > 0 {
			m.typeIdx--
		}
		return m, nil

	case "enter", "ctrl+l":
		if len(m.filteredTypes) == 0 {
			return m, nil
		}
		m.commitType = m.filteredTypes[m.typeIdx]
		m.typeInput.Blur()
		m.step = stepScope
		m.scopeInput.Focus()
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.typeInput, cmd = m.typeInput.Update(msg)
	// filter and reset cursor
	m.filteredTypes = filterScopes(commitTypes, m.typeInput.Value())
	m.typeIdx = 0

	return m, cmd
}

// filterScopes returns scopes that contain the query (case-insensitive).
func filterScopes(all []string, q string) []string {
	if q == "" {
		return all
	}
	q = strings.ToLower(q)
	var prefix, contains []string
	for _, s := range all {
		lower := strings.ToLower(s)
		if strings.HasPrefix(lower, q) {
			prefix = append(prefix, s)
		} else if strings.Contains(lower, q) {
			contains = append(contains, s)
		}
	}
	out := append(prefix, contains...)

	return out
}

// updateScope handles key events on the scope input step.
func (m model) updateScope(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc", "ctrl+h":
		// Go back to type step.
		m.scopeInput.Blur()
		m.scopeInput.SetValue("")
		m.filteredScopes = m.allScopes
		m.step = stepType
		return m, nil

	case "down", "ctrl+j":
		if m.scopeIdx < len(m.filteredScopes)-1 {
			m.scopeIdx++
			m.scopeInput.SetValue(m.filteredScopes[m.scopeIdx])
			m.scopeInput.CursorEnd()
		}
		return m, nil

	case "up", "ctrl+k":
		if m.scopeIdx > 0 {
			m.scopeIdx--
			m.scopeInput.SetValue(m.filteredScopes[m.scopeIdx])
			m.scopeInput.CursorEnd()
		}
		return m, nil

	case "enter", "ctrl+l":
		val := strings.TrimSpace(m.scopeInput.Value())
		if val == "" && len(m.filteredScopes) > 0 {
			val = m.filteredScopes[m.scopeIdx]
		}
		m.commitScope = val
		m.scopeInput.Blur()
		m.step = stepSubject
		m.subjectInput.Focus()
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.scopeInput, cmd = m.scopeInput.Update(msg)
	m.filteredScopes = filterScopes(m.allScopes, m.scopeInput.Value())
	m.scopeIdx = 0

	return m, cmd
}

// updateSubject handles key events on the subject input step.
func (m model) updateSubject(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc", "ctrl+h":
		// Go back to scope step.
		m.subjectInput.Blur()
		m.step = stepScope
		m.scopeInput.Focus()
		return m, textinput.Blink

	case "enter", "ctrl+l":
		subject := strings.TrimSpace(m.subjectInput.Value())
		if subject == "" {
			return m, nil // require non-empty subject
		}
		m.commitSubject = subject
		m.subjectInput.Blur()
		m.commitMsg = buildCommitMsg(m.commitType, m.commitScope, m.commitSubject)
		m.step = stepDone
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.subjectInput, cmd = m.subjectInput.Update(msg)

	full := buildCommitMsg(m.commitType, m.commitScope, m.subjectInput.Value())
	if len(full) > commitLimit {
		m.subjectInput.TextStyle = errorStyle
	} else {
		m.subjectInput.TextStyle = lipgloss.NewStyle()
	}

	return m, cmd
}

func buildCommitMsg(ctype, scope, subject string) string {
	base := ctype
	if scope != "" {
		base = fmt.Sprintf("%s(%s)", ctype, scope)
	}
	if subject == "" {
		return base
	}

	return fmt.Sprintf("%s: %s", base, subject)
}

func maxLineWidth(s string) int {
	max := 0
	for _, line := range strings.Split(s, "\n") {
		if w := lipgloss.Width(line); w > max {
			max = w
		}
	}
	return max
}

func stagedFilesView(files []StagedFile) string {
	if len(files) == 0 {
		return errorStyle.Render(" no staged files ")
	}
	addedStyle   := lipgloss.NewStyle().Foreground(colors.green)
	modifiedStyle := lipgloss.NewStyle().Foreground(colors.orange)
	deletedStyle := lipgloss.NewStyle().Foreground(colors.red)
	var lines []string
	for _, f := range files {
		switch f.Status {
		case "A":
			lines = append(lines, addedStyle.Render("+ "+f.Path))
		case "D":
			lines = append(lines, deletedStyle.Render("- "+f.Path))
		default:
			lines = append(lines, modifiedStyle.Render("~ "+f.Path))
		}
	}

	return strings.Join(lines, "\n")
}

func (m model) View() string {
	if m.step == stepDone {
		return borderStyle.Render(selectedStyle.Render(m.commitMsg)) + "\n"
	}

	var b strings.Builder

	switch m.step {
	case stepType:
		b.WriteString(titleStyle.Render(title) + "\n\n")
		b.WriteString(m.typeInput.View() + "\n")
		if len(m.filteredTypes) > 0 {
			b.WriteString("\n")
			for i, t := range m.filteredTypes {
				if i == m.typeIdx {
					b.WriteString(selectedStyle.Render(" 󰹹 "+t) + "\n")
				} else {
					b.WriteString(dimStyle.Render("   "+t) + "\n")
				}
			}
		}
		b.WriteString(previewLine(m))

	case stepScope:
		b.WriteString(titleStyle.Render(title) + "\n\n")
		b.WriteString(m.scopeInput.View() + "\n")
		if len(m.filteredScopes) > 0 {
			b.WriteString("\n")
			shown := m.filteredScopes
			if len(shown) > 5 {
				shown = shown[:5]
			}
			for i, s := range shown {
				if i == m.scopeIdx {
					b.WriteString(selectedStyle.Render(" 󰹹 "+s) + "\n")
				} else {
					b.WriteString(dimStyle.Render("   "+s) + "\n")
				}
			}
		}
		b.WriteString(previewLine(m))

	case stepSubject:
		b.WriteString(titleStyle.Render(title) + "\n\n")
		full := buildCommitMsg(m.commitType, m.commitScope, m.subjectInput.Value())
		count := fmt.Sprintf("%d/%d", len(full), commitLimit)
		var countStyle lipgloss.Style
		if len(full) > commitLimit {
			countStyle = errorStyle
		} else {
			countStyle = dimStyle
		}

		refWidth := lipgloss.Width(m.subjectInput.View())
		for _, candidate := range []int{
			lipgloss.Width(full),
			maxLineWidth(stagedFilesView(m.stagedFiles)),
			lipgloss.Width(titleStyle.Render(title)),
		} {
			if candidate > refWidth {
				refWidth = candidate
			}
		}

		b.WriteString(m.subjectInput.View() + "\n")
		padding := refWidth - lipgloss.Width(count)
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString(countStyle.Render(count) + "\n")
		b.WriteString(previewLine(m))
	}

	b.WriteString(stagedFilesView(m.stagedFiles))

	return borderStyle.Render(b.String()) + "\n"
}

func previewLine(m model) string {
	scope := strings.TrimSpace(m.scopeInput.Value())
	if m.step > stepScope {
		scope = m.commitScope
	}
	subject := strings.TrimSpace(m.subjectInput.Value())

	ctype := m.commitType
	if m.step == stepType && len(m.filteredTypes) > 0 {
		ctype = m.filteredTypes[m.typeIdx]
	}

	preview := buildCommitMsg(ctype, scope, subject)

	return "\n" + dimStyle.Render(preview) + "\n\n"
}
