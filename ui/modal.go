package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// overlayModal draws a dimmed full-screen background and composites the modal
// string centered on top (ANSI-aware).
func overlayModal(base, modal string, w, h int) string {
	if w <= 0 || h <= 0 {
		return base
	}

	dim := lipgloss.NewStyle().Faint(true).Render(base)
	lines := strings.Split(dim, "\n")
	padLinesToHeight(&lines, h)
	for i := range lines {
		lines[i] = lipgloss.PlaceHorizontal(w, lipgloss.Left, lines[i])
	}

	modalLines := strings.Split(modal, "\n")
	mw := lipgloss.Width(modal)
	mh := lipgloss.Height(modal)
	if mw <= 0 || mh <= 0 {
		return strings.Join(lines, "\n")
	}

	startRow := (h - mh) / 2
	startCol := (w - mw) / 2
	if startRow < 0 {
		startRow = 0
	}
	if startCol < 0 {
		startCol = 0
	}

	for i := 0; i < mh && startRow+i < len(lines); i++ {
		lines[startRow+i] = mergeAt(lines[startRow+i], modalLines[i], startCol, w)
	}

	return strings.Join(lines, "\n")
}

func mergeAt(bgLine, modalLine string, startCol, totalW int) string {
	left := ansi.Cut(bgLine, 0, startCol)
	mw := ansi.StringWidth(modalLine)
	right := ansi.Cut(bgLine, startCol+mw, totalW)
	return left + modalLine + right
}

func padLinesToHeight(lines *[]string, h int) {
	for len(*lines) < h {
		*lines = append(*lines, "")
	}
	if len(*lines) > h {
		*lines = (*lines)[:h]
	}
}
