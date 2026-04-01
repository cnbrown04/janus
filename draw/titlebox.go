package draw

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Border draws a rounded border with the title embedded in the top rule (╭─ Title ───╮).
// The output is exactly w cells wide and h lines tall when w,h >= 2.
func Border(title, body string, w, h int, border lipgloss.TerminalColor) string {
	if w < 2 || h < 2 {
		return ""
	}

	borderSt := lipgloss.NewStyle().Foreground(border)

	tl := borderSt.Render("╭")
	tr := borderSt.Render("╮")
	v := borderSt.Render("│")

	innerW := w - 2
	topMid := topBorderMiddle(innerW, title, border)
	topLine := tl + topMid + tr

	innerLines := h - 2
	bodyLines := splitBody(body, innerLines, innerW)

	var mid strings.Builder
	for i := 0; i < innerLines; i++ {
		var line string
		if i < len(bodyLines) {
			line = padOrTruncate(bodyLines[i], innerW)
		} else {
			line = strings.Repeat(" ", innerW)
		}
		mid.WriteString(v)
		mid.WriteString(line)
		mid.WriteString(v)
		if i < innerLines-1 {
			mid.WriteString("\n")
		}
	}

	bottomPlain := borderSt.Render("╰" + strings.Repeat("─", innerW) + "╯")

	if innerLines == 0 {
		return lipgloss.JoinVertical(lipgloss.Top, topLine, bottomPlain)
	}
	return lipgloss.JoinVertical(lipgloss.Top, topLine, mid.String(), bottomPlain)
}

func topBorderMiddle(innerW int, title string, border lipgloss.TerminalColor) string {
	borderSt := lipgloss.NewStyle().Foreground(border)
	hseg := "─"
	prefix := hseg + " "
	pw := ansi.StringWidth(prefix)
	spaceAfter := 1
	maxTitle := innerW - pw - spaceAfter
	if maxTitle < 0 {
		maxTitle = 0
	}
	t := title
	if ansi.StringWidth(t) > maxTitle {
		t = truncateToWidth(t, maxTitle)
	}
	tw := ansi.StringWidth(t)
	trailing := innerW - pw - tw - spaceAfter
	if trailing < 0 {
		trailing = 0
	}
	mid := prefix + t + strings.Repeat(" ", spaceAfter) + strings.Repeat(hseg, trailing)
	if ansi.StringWidth(mid) > innerW {
		mid = truncateToWidth(mid, innerW)
	} else if ansi.StringWidth(mid) < innerW {
		mid = mid + strings.Repeat(hseg, innerW-ansi.StringWidth(mid))
	}
	return borderSt.Render(mid)
}

func splitBody(body string, maxLines, lineWidth int) []string {
	if body == "" {
		return nil
	}
	raw := strings.Split(body, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		if len(out) >= maxLines {
			break
		}
		out = append(out, truncateToWidth(line, lineWidth))
	}
	return out
}

func padOrTruncate(s string, w int) string {
	sw := ansi.StringWidth(s)
	if sw > w {
		return truncateToWidth(s, w)
	}
	if sw < w {
		return s + strings.Repeat(" ", w-sw)
	}
	return s
}

func truncateToWidth(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if ansi.StringWidth(s) <= max {
		return s
	}
	var b strings.Builder
	w := 0
	for _, r := range s {
		rs := string(r)
		rw := ansi.StringWidth(rs)
		if w+rw > max {
			break
		}
		b.WriteString(rs)
		w += rw
	}
	return b.String()
}
