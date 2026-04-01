package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/cnbrown04/janus/draw"
)

// renderQueryBody draws the query buffer with SQL highlighting, scroll, and cursor
// aligned to textarea.LineInfo (same wrap algorithm as bubbles/textarea).
func renderQueryBody(m *Model, innerW, innerH int) string {
	m.queryArea.SetWidth(innerW)
	m.queryArea.SetHeight(innerH)

	val := m.queryArea.Value()
	if val == "" && m.queryArea.Line() == 0 && m.queryArea.Placeholder != "" {
		return padQueryBlock(m.queryArea.View(), innerW, innerH)
	}

	selLo, selHi := -1, -1
	if m.querySelAnchorLine >= 0 {
		a, c := m.querySelAnchorLine, m.queryArea.Line()
		selLo = min(a, c)
		selHi = max(a, c)
	}

	return renderQueryHighlightedView(m, innerW, innerH, val, selLo, selHi)
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

func cursorDisplayRow(m *Model) int {
	lines := strings.Split(m.queryArea.Value(), "\n")
	w := m.queryArea.Width()
	if w < 1 {
		w = 1
	}
	row := 0
	curLn := m.queryArea.Line()
	for i := 0; i < curLn && i < len(lines); i++ {
		row += len(wrapRunes([]rune(lines[i]), w))
	}
	row += m.queryArea.LineInfo().RowOffset
	return row
}

func buildHighlightedDisplayLines(val string, innerW int) []string {
	if innerW < 1 {
		innerW = 1
	}
	logical := strings.Split(val, "\n")
	var out []string
	for _, ln := range logical {
		for _, seg := range wrapRunes([]rune(ln), innerW) {
			out = append(out, highlightSQLLine(string(seg)))
		}
	}
	return out
}

func (m *Model) syncQueryScroll(innerH, totalLines int) {
	if innerH < 1 {
		innerH = 1
	}
	cur := cursorDisplayRow(m)
	if cur < m.queryScrollOff {
		m.queryScrollOff = cur
	}
	if cur >= m.queryScrollOff+innerH {
		m.queryScrollOff = cur - innerH + 1
	}
	maxOff := max(0, totalLines-innerH)
	if m.queryScrollOff > maxOff {
		m.queryScrollOff = maxOff
	}
}

func mapDisplayRowToLogicalLine(displayRow int, logicalLines []string, w int) int {
	r := 0
	for li, ln := range logicalLines {
		n := len(wrapRunes([]rune(ln), w))
		if n < 1 {
			n = 1
		}
		if displayRow < r+n {
			return li
		}
		r += n
	}
	if len(logicalLines) == 0 {
		return 0
	}
	return len(logicalLines) - 1
}

func renderQueryHighlightedView(m *Model, innerW, innerH int, val string, selLo, selHi int) string {
	if val == "" {
		m.syncQueryScroll(innerH, 1)
		line := strings.Repeat(" ", innerW)
		if m.queryArea.Focused() {
			line = injectWideCursor(m, line, innerW, m.queryArea.LineInfo().CharOffset)
		}
		return padQueryBlock(line, innerW, innerH)
	}

	display := buildHighlightedDisplayLines(val, innerW)
	m.syncQueryScroll(innerH, len(display))

	end := min(m.queryScrollOff+innerH, len(display))
	window := display[m.queryScrollOff:end]

	selSt := lipgloss.NewStyle().Background(draw.SelectionBG).Foreground(draw.SelectionFG)
	curRow := cursorDisplayRow(m)
	logicalLines := strings.Split(val, "\n")
	w := m.queryArea.Width()
	if w < 1 {
		w = 1
	}

	var rowStrs []string
	for i, hl := range window {
		globalRow := m.queryScrollOff + i
		line := padANSIToWidth(hl, innerW)

		if selLo >= 0 && selHi >= selLo {
			if lr := mapDisplayRowToLogicalLine(globalRow, logicalLines, w); lr >= selLo && lr <= selHi {
				line = selSt.Render(line)
			}
		}

		if m.queryArea.Focused() && globalRow == curRow {
			col := m.queryArea.LineInfo().CharOffset
			line = injectWideCursor(m, line, innerW, col)
		}

		rowStrs = append(rowStrs, line)
	}

	return padQueryBlock(strings.Join(rowStrs, "\n"), innerW, innerH)
}

func injectWideCursor(m *Model, line string, innerW, col int) string {
	sw := ansi.StringWidth(line)
	if col < 0 {
		col = 0
	}
	if col > sw {
		line = line + strings.Repeat(" ", col-sw)
	}

	mid := ansi.Cut(line, col, col+1)
	if mid == "" {
		m.queryArea.Cursor.SetChar(" ")
	} else {
		m.queryArea.Cursor.SetChar(mid)
	}
	if m.queryArea.Focused() {
		st := m.queryArea.FocusedStyle
		m.queryArea.Cursor.TextStyle = st.CursorLine.Inherit(st.Base).Inline(true)
	} else {
		st := m.queryArea.BlurredStyle
		m.queryArea.Cursor.TextStyle = st.Text.Inherit(st.Base).Inline(true)
	}

	rightEdge := ansi.StringWidth(line)
	if col+1 > rightEdge {
		rightEdge = col + 1
	}
	left := ansi.Cut(line, 0, col)
	right := ansi.Cut(line, col+1, rightEdge)
	return left + m.queryArea.Cursor.View() + right
}
