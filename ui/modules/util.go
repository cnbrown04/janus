package modules

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// PadLineToWidth pads or truncates a plain or ANSI string to maxW terminal cells.
func PadLineToWidth(s string, maxW int) string {
	sw := ansi.StringWidth(s)
	if sw > maxW {
		return ansi.Truncate(s, maxW, "")
	}
	if sw < maxW {
		return s + strings.Repeat(" ", maxW-sw)
	}
	return s
}

// PadQueryBlock pads and clips a multi-line block to w×h terminal cells (ANSI-aware per line).
func PadQueryBlock(s string, w, h int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	for i := range lines {
		lines[i] = PadLineToWidth(lines[i], w)
	}
	for len(lines) < h {
		lines = append(lines, strings.Repeat(" ", w))
	}
	return strings.Join(lines, "\n")
}

func padANSIToWidth(s string, w int) string {
	sw := ansi.StringWidth(s)
	if sw > w {
		return ansi.Truncate(s, w, "")
	}
	if sw < w {
		return s + strings.Repeat(" ", w-sw)
	}
	return s
}
