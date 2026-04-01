package ui

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// padQueryBlock pads and clips a multi-line block to w×h terminal cells (ANSI-aware per line).
func padQueryBlock(s string, w, h int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for i := range lines {
		lines[i] = padLineToWidth(lines[i], w)
	}
	for len(lines) < h {
		lines = append(lines, strings.Repeat(" ", w))
	}
	return strings.Join(lines, "\n")
}

func padLineToWidth(s string, maxW int) string {
	sw := ansi.StringWidth(s)
	if sw > maxW {
		return ansi.Truncate(s, maxW, "")
	}
	if sw < maxW {
		return s + strings.Repeat(" ", maxW-sw)
	}
	return s
}
