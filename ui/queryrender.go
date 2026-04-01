package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/cnbrown04/janus/draw"
)

// renderQueryBody draws the query buffer: full textarea view, or line-highlighted text when a line range is active.
func renderQueryBody(m *Model, innerW, innerH int) string {
	m.queryArea.SetWidth(innerW)
	m.queryArea.SetHeight(innerH)

	anchor := m.querySelAnchorLine
	cur := m.queryArea.Line()
	hasAnchor := anchor >= 0
	lo, hi := 0, 0
	if hasAnchor {
		lo = min(anchor, cur)
		hi = max(anchor, cur)
	}

	if !hasAnchor {
		return padQueryBlock(m.queryArea.View(), innerW, innerH)
	}

	return renderQueryLineSelection(m, innerW, innerH, lo, hi)
}

func renderQueryLineSelection(m *Model, innerW, innerH int, lo, hi int) string {
	val := m.queryArea.Value()
	lines := strings.Split(val, "\n")
	selSt := lipgloss.NewStyle().Background(draw.SelectionBG).Foreground(draw.SelectionFG)

	var b strings.Builder
	for i, ln := range lines {
		if i > 0 {
			b.WriteString("\n")
		}
		line := padLineToWidth(ln, innerW)
		if i >= lo && i <= hi {
			b.WriteString(selSt.Render(line))
		} else {
			b.WriteString(line)
		}
	}
	return padQueryBlock(b.String(), innerW, innerH)
}
