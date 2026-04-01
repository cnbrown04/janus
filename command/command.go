package command

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/cnbrown04/janus/draw"
)

// Model is the ex-style command line (":" to open). Text entry will later drive nvim-like commands.
type Model struct {
	input textinput.Model
	open  bool
}

// New builds a command-line model with default text input settings.
func New() Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 2048
	ti.Placeholder = ""
	fg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	ti.TextStyle = fg
	ti.Cursor.Style = fg
	ti.Cursor.TextStyle = fg
	return Model{input: ti}
}

func (m Model) Active() bool {
	return m.open
}

// Toggle opens or closes the command line. When opened, the field is focused and blink starts.
func (m Model) Toggle() (Model, tea.Cmd) {
	m.open = !m.open
	if m.open {
		m.input.Focus()
		return m, textinput.Blink
	}
	m.input.Blur()
	return m, nil
}

// Close hides the command line and blurs input.
func (m Model) Close() Model {
	m.open = false
	m.input.Blur()
	return m
}

// Update routes messages to the text field while the command line is active.
// Enter runs the current line via exec then closes the bar.
func Update(m Model, msg tea.Msg) (Model, tea.Cmd) {
	if !m.open {
		return m, nil
	}

	if km, ok := msg.(tea.KeyMsg); ok && km.Type == tea.KeyEnter {
		line := strings.TrimSpace(m.input.Value())
		m.input.SetValue("")
		m.open = false
		m.input.Blur()
		return m, exec(line)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// exec runs a single ex-style command. Only q and quit exit the program for now.
func exec(line string) tea.Cmd {
	switch strings.ToLower(line) {
	case "q", "quit":
		return tea.Quit
	default:
		return nil
	}
}

// View renders a one-row command bar inside a titled border. termWidth bounds the box width.
func (m Model) View(termWidth int) string {
	if !m.open || termWidth <= 0 {
		return ""
	}

	bw := min(56, termWidth-4)
	if bw < 16 {
		bw = max(12, termWidth-2)
	}
	if bw < 4 {
		bw = termWidth
	}

	innerW := bw - 2
	if innerW < 1 {
		innerW = 1
	}

	inp := m.input
	inp.Width = innerW
	body := padOrTruncate(inp.View(), innerW)

	return draw.Border("Command", body, bw, 3, lipgloss.Color("#FFFFFF"))
}

func padOrTruncate(s string, w int) string {
	sw := ansi.StringWidth(s)
	if sw > w {
		return ansi.Truncate(s, w, "")
	}
	if sw < w {
		return s + strings.Repeat(" ", w-sw)
	}
	return s
}
