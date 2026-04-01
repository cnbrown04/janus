package modules

import (
	"strings"
	"sync"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	rw "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/draw"
)

// QueryRenderContext holds query-editor state for rendering (textarea + scroll + line selection).
type QueryRenderContext struct {
	TA            *textarea.Model
	SelAnchorLine int
	ScrollOff     *int
}

// --- sql wrap (from bubbles/textarea, MIT) ---

func repeatSpaces(n int) []rune {
	if n <= 0 {
		return nil
	}
	return []rune(strings.Repeat(" ", n))
}

func wrapRunes(runes []rune, width int) [][]rune {
	if width < 1 {
		width = 1
	}
	if len(runes) == 0 {
		return [][]rune{{}}
	}
	var (
		lines  = [][]rune{{}}
		word   = []rune{}
		row    int
		spaces int
	)

	for _, r := range runes {
		if unicode.IsSpace(r) {
			spaces++
		} else {
			word = append(word, r)
		}

		if spaces > 0 {
			if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces > width {
				row++
				lines = append(lines, []rune{})
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			} else {
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			}
		} else {
			if len(word) == 0 {
				continue
			}
			lastCharLen := rw.RuneWidth(word[len(word)-1])
			if uniseg.StringWidth(string(word))+lastCharLen > width {
				if len(lines[row]) > 0 {
					row++
					lines = append(lines, []rune{})
				}
				lines[row] = append(lines[row], word...)
				word = nil
			}
		}
	}

	if uniseg.StringWidth(string(lines[row]))+uniseg.StringWidth(string(word))+spaces >= width {
		lines = append(lines, []rune{})
		lines[row+1] = append(lines[row+1], word...)
		spaces++
		lines[row+1] = append(lines[row+1], repeatSpaces(spaces)...)
	} else {
		lines[row] = append(lines[row], word...)
		spaces++
		lines[row] = append(lines[row], repeatSpaces(spaces)...)
	}

	return lines
}

// --- sql highlight ---

var (
	sqlComment  = lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")).Italic(true)
	sqlKeyword  = lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934"))
	sqlString   = lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26"))
	sqlNumber   = lipgloss.NewStyle().Foreground(lipgloss.Color("#d3869b"))
	sqlName     = lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598"))
	sqlOperator = lipgloss.NewStyle().Foreground(lipgloss.Color("#fe8019"))
	sqlPunct    = lipgloss.NewStyle().Foreground(lipgloss.Color("#a89984"))
	sqlDefault  = lipgloss.NewStyle()
)

var (
	sqlLexerOnce sync.Once
	sqlLexer     chroma.Lexer
)

func initSQLLexer() {
	sqlLexerOnce.Do(func() {
		l := lexers.Get("sql")
		if l == nil {
			l = lexers.Fallback
		}
		sqlLexer = chroma.Coalesce(l)
	})
}

func styleForSQLToken(t chroma.TokenType) lipgloss.Style {
	switch t {
	case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration,
		chroma.KeywordPseudo, chroma.KeywordReserved, chroma.KeywordType:
		return sqlKeyword
	case chroma.String, chroma.StringChar, chroma.StringDouble, chroma.StringBacktick,
		chroma.StringHeredoc, chroma.StringSingle, chroma.StringInterpol, chroma.StringAffix:
		return sqlString
	case chroma.Number, chroma.NumberInteger, chroma.NumberIntegerLong, chroma.NumberFloat, chroma.NumberHex,
		chroma.NumberOct, chroma.NumberBin:
		return sqlNumber
	case chroma.Name, chroma.NameAttribute, chroma.NameBuiltin, chroma.NameClass, chroma.NameDecorator,
		chroma.NameEntity, chroma.NameException, chroma.NameFunction, chroma.NameProperty, chroma.NameLabel,
		chroma.NameNamespace, chroma.NameOther, chroma.NameTag, chroma.NameVariable,
		chroma.NameConstant, chroma.NameBuiltinPseudo, chroma.NameFunctionMagic:
		return sqlName
	case chroma.Comment, chroma.CommentSingle, chroma.CommentMultiline, chroma.CommentPreproc,
		chroma.CommentSpecial, chroma.CommentHashbang:
		return sqlComment
	case chroma.Operator, chroma.OperatorWord:
		return sqlOperator
	case chroma.Punctuation:
		return sqlPunct
	default:
		return sqlDefault
	}
}

func highlightSQLLine(line string) string {
	initSQLLexer()
	it, err := sqlLexer.Tokenise(nil, line)
	if err != nil {
		return line
	}
	var b strings.Builder
	for t := it(); t != chroma.EOF; t = it() {
		if t.Value == "" {
			continue
		}
		b.WriteString(styleForSQLToken(t.Type).Render(t.Value))
	}
	return b.String()
}

// RenderQueryBody draws the query buffer with SQL highlighting, scroll, and cursor.
func RenderQueryBody(ctx *QueryRenderContext, innerW, innerH int) string {
	ta := ctx.TA
	ta.SetWidth(innerW)
	ta.SetHeight(innerH)

	val := ta.Value()
	if val == "" && ta.Line() == 0 && ta.Placeholder != "" {
		return PadQueryBlock(ta.View(), innerW, innerH)
	}

	selLo, selHi := -1, -1
	if ctx.SelAnchorLine >= 0 {
		a, c := ctx.SelAnchorLine, ta.Line()
		selLo = min(a, c)
		selHi = max(a, c)
	}

	return renderQueryHighlightedView(ctx, innerW, innerH, val, selLo, selHi)
}

func cursorDisplayRow(ta *textarea.Model) int {
	lines := strings.Split(ta.Value(), "\n")
	w := ta.Width()
	if w < 1 {
		w = 1
	}
	row := 0
	curLn := ta.Line()
	for i := 0; i < curLn && i < len(lines); i++ {
		row += len(wrapRunes([]rune(lines[i]), w))
	}
	row += ta.LineInfo().RowOffset
	return row
}

func buildDisplayLinePairs(val string, innerW int) (plain []string, syntax []string) {
	if innerW < 1 {
		innerW = 1
	}
	logical := strings.Split(val, "\n")
	for _, ln := range logical {
		for _, seg := range wrapRunes([]rune(ln), innerW) {
			ps := string(seg)
			plain = append(plain, ps)
			syntax = append(syntax, highlightSQLLine(ps))
		}
	}
	return plain, syntax
}

func syncQueryScroll(ta *textarea.Model, scrollOff *int, innerH, totalLines int) {
	if innerH < 1 {
		innerH = 1
	}
	cur := cursorDisplayRow(ta)
	if cur < *scrollOff {
		*scrollOff = cur
	}
	if cur >= *scrollOff+innerH {
		*scrollOff = cur - innerH + 1
	}
	maxOff := max(0, totalLines-innerH)
	if *scrollOff > maxOff {
		*scrollOff = maxOff
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

func renderQueryHighlightedView(ctx *QueryRenderContext, innerW, innerH int, val string, selLo, selHi int) string {
	ta := ctx.TA
	scrollOff := ctx.ScrollOff
	if val == "" {
		syncQueryScroll(ta, scrollOff, innerH, 1)
		line := strings.Repeat(" ", innerW)
		if ta.Focused() {
			line = injectWideCursor(ta, line, innerW, ta.LineInfo().CharOffset)
		}
		return PadQueryBlock(line, innerW, innerH)
	}

	selSt := lipgloss.NewStyle().Background(draw.SelectionBG).Foreground(draw.SelectionFG)
	curRow := cursorDisplayRow(ta)
	logicalLines := strings.Split(val, "\n")
	w := ta.Width()
	if w < 1 {
		w = 1
	}

	plainLines, synLines := buildDisplayLinePairs(val, innerW)
	syncQueryScroll(ta, scrollOff, innerH, len(plainLines))

	end := min(*scrollOff+innerH, len(plainLines))
	windowPlain := plainLines[*scrollOff:end]
	windowSyn := synLines[*scrollOff:end]

	var rowStrs []string
	for i := range windowPlain {
		globalRow := *scrollOff + i

		var line string
		if selLo >= 0 && selHi >= selLo {
			lr := mapDisplayRowToLogicalLine(globalRow, logicalLines, w)
			if lr >= selLo && lr <= selHi {
				line = selSt.Render(PadLineToWidth(windowPlain[i], innerW))
			} else {
				line = padANSIToWidth(windowSyn[i], innerW)
			}
		} else {
			line = padANSIToWidth(windowSyn[i], innerW)
		}

		if ta.Focused() && globalRow == curRow {
			col := ta.LineInfo().CharOffset
			line = injectWideCursor(ta, line, innerW, col)
		}

		rowStrs = append(rowStrs, line)
	}

	return PadQueryBlock(strings.Join(rowStrs, "\n"), innerW, innerH)
}

func injectWideCursor(ta *textarea.Model, line string, innerW, col int) string {
	sw := ansi.StringWidth(line)
	if col < 0 {
		col = 0
	}
	if col > sw {
		line = line + strings.Repeat(" ", col-sw)
	}

	mid := ansi.Cut(line, col, col+1)
	if mid == "" {
		ta.Cursor.SetChar(" ")
	} else {
		ta.Cursor.SetChar(mid)
	}
	if ta.Focused() {
		st := ta.FocusedStyle
		ta.Cursor.TextStyle = st.CursorLine.Inherit(st.Base).Inline(true)
	} else {
		st := ta.BlurredStyle
		ta.Cursor.TextStyle = st.Text.Inherit(st.Base).Inline(true)
	}

	rightEdge := ansi.StringWidth(line)
	if col+1 > rightEdge {
		rightEdge = col + 1
	}
	left := ansi.Cut(line, 0, col)
	right := ansi.Cut(line, col+1, rightEdge)
	return left + ta.Cursor.View() + right
}

// MatchesQueryExecute reports keys that should run the query without inserting text.
func MatchesQueryExecute(msg tea.KeyMsg, keys bindings.KeyMap) bool {
	if key.Matches(msg, keys.QueryExecute) {
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

// QueryExecuteText returns the SQL to run from a line-range selection or full buffer.
func QueryExecuteText(val string, anchorLine, curLine int) string {
	lines := strings.Split(val, "\n")
	if anchorLine >= 0 {
		lo, hi := min(anchorLine, curLine), max(anchorLine, curLine)
		if len(lines) > 0 && lo <= hi {
			lo = min(max(lo, 0), len(lines)-1)
			hi = min(max(hi, 0), len(lines)-1)
			return strings.Join(lines[lo:hi+1], "\n")
		}
	}
	return val
}

// ShouldForwardQueryNormalMode is true for navigation keys passed to the textarea in normal mode.
func ShouldForwardQueryNormalMode(msg tea.KeyMsg) bool {
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
		return false
	case tea.KeyTab, tea.KeyShiftTab:
		return false
	}
	return true
}

// RemapShiftVerticalForTextarea maps shift+vertical arrows to plain arrows for the textarea.
func RemapShiftVerticalForTextarea(msg tea.KeyMsg) tea.Msg {
	switch msg.Type {
	case tea.KeyShiftUp:
		return tea.KeyMsg{Type: tea.KeyUp}
	case tea.KeyShiftDown:
		return tea.KeyMsg{Type: tea.KeyDown}
	default:
		return msg
	}
}
