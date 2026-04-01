package ui

import tea "github.com/charmbracelet/bubbletea"

// shouldForwardQueryNormalMode is true for navigation keys we pass to the textarea in normal (non-insert) mode.
// Typed characters, space, tab, enter, and paste stay in normal mode (no insertion).
func shouldForwardQueryNormalMode(msg tea.KeyMsg) bool {
	if msg.Paste {
		return false
	}
	switch msg.Type {
	case tea.KeyRunes:
		if len(msg.Runes) > 0 {
			return false
		}
	case tea.KeySpace:
		return false
	case tea.KeyEnter:
		// KeyCtrlM is the same KeyType as KeyEnter (carriage return).
		return false
	case tea.KeyTab, tea.KeyShiftTab:
		return false
	}
	return true
}

// remapShiftVerticalForTextarea turns shift+arrow into plain vertical arrows so the
// textarea moves the cursor (line selection uses the anchor + current line).
func remapShiftVerticalForTextarea(msg tea.KeyMsg) tea.Msg {
	switch msg.Type {
	case tea.KeyShiftUp:
		return tea.KeyMsg{Type: tea.KeyUp}
	case tea.KeyShiftDown:
		return tea.KeyMsg{Type: tea.KeyDown}
	default:
		return msg
	}
}
