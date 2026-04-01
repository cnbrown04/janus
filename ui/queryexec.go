package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// matchesQueryExecute reports keys that should run the query without inserting text.
// Note: ctrl+enter is usually indistinguishable from enter in terminals (both are carriage return).
func (m *Model) matchesQueryExecute(msg tea.KeyMsg) bool {
	if key.Matches(msg, m.keys.QueryExecute) {
		return true
	}
	if msg.Type == tea.KeyCtrlJ {
		return true
	}
	if msg.Type == tea.KeyEnter && msg.Alt {
		return true
	}
	return false
}

// runQueryExecute runs the query: selected logical lines if a line range is active, otherwise the full buffer.
func (m *Model) runQueryExecute() {
	val := m.queryArea.Value()
	lines := strings.Split(val, "\n")
	anchor := m.querySelAnchorLine
	cur := m.queryArea.Line()

	var text string
	if anchor >= 0 {
		lo, hi := min(anchor, cur), max(anchor, cur)
		if len(lines) > 0 && lo <= hi {
			lo = min(max(lo, 0), len(lines)-1)
			hi = min(max(hi, 0), len(lines)-1)
			text = strings.Join(lines[lo:hi+1], "\n")
		} else {
			text = val
		}
	} else {
		text = val
	}

	_ = text // TODO: execute against connection / results panel
	m.querySelAnchorLine = -1
}
