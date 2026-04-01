package modules

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/draw"
)

// DatabasePanelState holds UI state for the Database panel (dropdown + current selection).
type DatabasePanelState struct {
	DropdownOpen bool
	Cursor       int
	Options      []string
	// Selected is always shown when the dropdown is closed (current connection).
	Selected string
}

// DefaultDatabaseOptions returns stub connection labels.
func DefaultDatabaseOptions() []string {
	return []string{
		"postgres (local)",
		"postgres (staging)",
		"sqlite",
		"mysql",
		"dev (mock)",
	}
}

func (s *DatabasePanelState) ensureOptions() {
	if len(s.Options) == 0 {
		s.Options = DefaultDatabaseOptions()
	}
	if s.Selected == "" && len(s.Options) > 0 {
		s.Selected = s.Options[0]
	}
}

// RenderDatabasePanel draws the Database panel; focused highlights the border when databaseFocused.
func RenderDatabasePanel(w, h int, st *DatabasePanelState, databaseFocused bool) string {
	st.ensureOptions()
	innerW := w - 2
	innerH := h - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}
	var body string
	if st.DropdownOpen && len(st.Options) > 0 {
		maxLines := min(len(st.Options), innerH)
		start := 0
		if st.Cursor >= maxLines {
			start = st.Cursor - maxLines + 1
		}
		var lines []string
		for i := 0; i < maxLines; i++ {
			idx := start + i
			label := st.Options[idx]
			prefix := "  "
			if idx == st.Cursor {
				prefix = "> "
			}
			lines = append(lines, PadLineToWidth(prefix+label, innerW))
		}
		body = strings.Join(lines, "\n")
	} else {
		line := st.Selected
		if line == "" {
			line = st.Options[0]
		}
		body = PadLineToWidth(line, innerW)
	}
	return draw.Border("Database", body, w, h, draw.PanelBorder(databaseFocused))
}

// HandleDatabaseKey handles keys when the Database panel is focused. Returns true if consumed.
func HandleDatabaseKey(msg tea.KeyMsg, keys bindings.KeyMap, st *DatabasePanelState) bool {
	st.ensureOptions()

	if st.DropdownOpen && key.Matches(msg, keys.Close) {
		st.DropdownOpen = false
		return true
	}

	if msg.Type == tea.KeyEnter {
		if !st.DropdownOpen {
			st.DropdownOpen = true
			st.Cursor = indexOf(st.Options, st.Selected)
			if st.Cursor < 0 {
				st.Cursor = 0
			}
		} else {
			if len(st.Options) > 0 && st.Cursor >= 0 && st.Cursor < len(st.Options) {
				st.Selected = st.Options[st.Cursor]
			}
			st.DropdownOpen = false
		}
		return true
	}

	if !st.DropdownOpen {
		return false
	}

	switch {
	case key.Matches(msg, keys.Up), msg.Type == tea.KeyUp:
		if st.Cursor > 0 {
			st.Cursor--
		}
		return true
	case key.Matches(msg, keys.Down), msg.Type == tea.KeyDown:
		if st.Cursor < len(st.Options)-1 {
			st.Cursor++
		}
		return true
	default:
		return false
	}
}

func indexOf(opts []string, s string) int {
	for i, o := range opts {
		if o == s {
			return i
		}
	}
	return -1
}
