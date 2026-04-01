package draw

import "github.com/charmbracelet/lipgloss"

// PanelBorder returns a lipgloss border color: focus green when the panel is selected, otherwise muted.
func PanelBorder(selected bool) lipgloss.TerminalColor {
	if selected {
		return BorderFocus
	}
	return BorderMuted
}
