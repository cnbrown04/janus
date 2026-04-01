package draw

import "github.com/charmbracelet/lipgloss"

// BorderMuted is the default border color for unfocused panels.
var BorderMuted = lipgloss.Color("245")

// BorderFocus is the border color for the active panel (muted green, nvim-adjacent).
var BorderFocus = lipgloss.Color("#689d6a")

// SelectionBG / SelectionFG are used for full-line visual selection in the query buffer.
var (
	SelectionBG = lipgloss.Color("#3d5a40")
	SelectionFG = lipgloss.Color("#d5ecd8")
)
