package modules

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/evertras/bubble-table/table"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/draw"
)

// colPadKey is a synthetic flex column that absorbs leftover width (empty cells).
const colPadKey = "__janus_pad__"

type resultsTableKind int

const (
	resultsKindUser resultsTableKind = iota
	resultsKindTall
	resultsKindWide
)

// ResultsTable is an interactive bubble-table + cell editor for the Results panel (test / stub).
type ResultsTable struct {
	active bool
	kind   resultsTableKind

	colKeys   []string
	colTitles []string
	data      [][]string
	tbl       table.Model

	col     int
	editing bool
	editor  textinput.Model
}

var cellFocusStyle = lipgloss.NewStyle().
	Background(draw.SelectionBG).
	Foreground(draw.SelectionFG)

var rowHiStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#3a3a3a")).
	Foreground(lipgloss.Color("252"))

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func defaultResultsKeyMap() table.KeyMap {
	km := table.DefaultKeyMap()
	km.RowSelectToggle = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "select (unused)"),
	)
	km.Filter = key.NewBinding(
		key.WithKeys("f8"),
		key.WithHelp("f8", "filter (off)"),
	)
	return km
}

// wideResultsKeyMap adds reliable column-scroll keys. Many terminals do not send shift+arrow as shift+left/right.
func wideResultsKeyMap() table.KeyMap {
	km := defaultResultsKeyMap()
	km.ScrollLeft = key.NewBinding(
		key.WithKeys("shift+left", "["),
		key.WithHelp("[ / shift+←", "scroll cols ←"),
	)
	km.ScrollRight = key.NewBinding(
		key.WithKeys("shift+right", "]"),
		key.WithHelp("] / shift+→", "scroll cols →"),
	)
	return km
}

func baseTableShell(cols []table.Column, km table.KeyMap) table.Model {
	return table.New(cols).
		Filtered(false).
		WithKeyMap(km).
		HighlightStyle(rowHiStyle)
}

// buildMeasuredColumns returns fixed-width columns sized to header + cell content (ANSI-aware).
// If minColWidth > 0, each column is at least that wide (characters) so wide demos overflow typical panels.
func buildMeasuredColumns(keys, titles []string, data [][]string, minColWidth int) []table.Column {
	cols := make([]table.Column, len(keys))
	for j, k := range keys {
		title := k
		if j < len(titles) && titles[j] != "" {
			title = titles[j]
		}
		w := ansi.StringWidth(title)
		for _, row := range data {
			if j < len(row) {
				cw := ansi.StringWidth(row[j])
				if cw > w {
					w = cw
				}
			}
		}
		if w < 2 {
			w = 2
		}
		cellW := w + 2
		if minColWidth > 0 && cellW < minColWidth {
			cellW = minColWidth
		}
		cols[j] = table.NewColumn(k, title, cellW)
	}
	return cols
}

func finishResultsTable(kind resultsTableKind, colKeys, colTitles []string, data [][]string) ResultsTable {
	minW := 0
	if kind == resultsKindWide {
		minW = 11
	}
	dc := buildMeasuredColumns(colKeys, colTitles, data, minW)
	ed := textinput.New()
	ed.CharLimit = 256
	ed.Prompt = ""

	var tbl table.Model
	switch kind {
	case resultsKindWide:
		tbl = baseTableShell(slices.Clone(dc), wideResultsKeyMap())
	default:
		flex := table.NewFlexColumn(colPadKey, " ", 1)
		tbl = baseTableShell(append(slices.Clone(dc), flex), defaultResultsKeyMap())
	}

	rt := ResultsTable{
		active:    true,
		kind:      kind,
		colKeys:   colKeys,
		colTitles: colTitles,
		data:      data,
		tbl:       tbl,
		editor:    ed,
	}
	rt.applyDimensions(80, 24)
	rt.tbl = rt.tbl.WithRows(rt.buildRowsForHi(0)).WithHighlightedRow(0)
	return rt
}

// NewResultsTableForCatalog opens a demo table by catalog name, or inactive ResultsTable if unknown.
func NewResultsTableForCatalog(name string) ResultsTable {
	switch name {
	case "user":
		return NewResultsUserTable()
	case "tall":
		return NewResultsTallTable()
	case "wide":
		return NewResultsWideTable()
	default:
		return ResultsTable{}
	}
}

// NewResultsUserTable is a small editable user grid with a trailing flex padding column.
func NewResultsUserTable() ResultsTable {
	colKeys := []string{"id", "email", "name", "role", "created"}
	colTitles := []string{"id", "email", "name", "role", "created"}
	data := [][]string{
		{"1", "alex.morgan@example.com", "Alex Morgan", "admin", "2024-06-12"},
		{"2", "jordan.lee@example.com", "Jordan Lee", "editor", "2025-01-03"},
		{"3", "sam.taylor@example.net", "Sam Taylor", "viewer", "2025-08-22"},
		{"4", "riley.chen@example.com", "Riley Chen", "editor", "2026-02-14"},
		{"5", "casey.walker@example.org", "Casey Walker", "viewer", "2026-03-01"},
	}
	return finishResultsTable(resultsKindUser, colKeys, colTitles, data)
}

// NewResultsTallTable is many rows with pagination (vertical overflow).
func NewResultsTallTable() ResultsTable {
	keys, titles, data := newTallDemoTable()
	return finishResultsTable(resultsKindTall, keys, titles, data)
}

// NewResultsWideTable is many columns with horizontal scrolling (horizontal overflow).
func NewResultsWideTable() ResultsTable {
	keys, titles, data := newWideDemoTable()
	return finishResultsTable(resultsKindWide, keys, titles, data)
}

func (rt *ResultsTable) applyDimensions(innerW, innerH int) {
	if innerW < 12 {
		innerW = 12
	}
	if innerH < 4 {
		innerH = 4
	}
	minW := 0
	if rt.kind == resultsKindWide {
		minW = 11
	}
	dc := buildMeasuredColumns(rt.colKeys, rt.colTitles, rt.data, minW)
	switch rt.kind {
	case resultsKindWide:
		rt.tbl = rt.tbl.WithColumns(slices.Clone(dc)).
			WithTargetWidth(0).
			WithMaxTotalWidth(innerW).
			WithNoPagination().
			WithFooterVisibility(false).
			WithHorizontalFreezeColumnCount(1)
	case resultsKindTall:
		flex := table.NewFlexColumn(colPadKey, " ", 1)
		cols := append(slices.Clone(dc), flex)
		ps := clampInt(innerH-9, 5, 200)
		rt.tbl = rt.tbl.WithColumns(cols).
			WithTargetWidth(innerW).
			WithMaxTotalWidth(innerW).
			WithPageSize(ps).
			WithFooterVisibility(true)
	case resultsKindUser:
		flex := table.NewFlexColumn(colPadKey, " ", 1)
		cols := append(slices.Clone(dc), flex)
		rt.tbl = rt.tbl.WithColumns(cols).
			WithTargetWidth(innerW).
			WithMaxTotalWidth(innerW).
			WithNoPagination().
			WithFooterVisibility(false).
			WithHorizontalFreezeColumnCount(0)
	}
}

// Active reports whether the Results panel should render the interactive table.
func (rt *ResultsTable) Active() bool {
	return rt != nil && rt.active
}

// Blur commits any in-progress edit and marks the inner table unfocused for updates.
func (rt *ResultsTable) Blur() {
	if rt == nil || !rt.active {
		return
	}
	if rt.editing {
		rt.commitEdit()
	}
	rt.tbl = rt.tbl.Focused(false)
}

func resultsKeyPassThrough(msg tea.KeyMsg, keys bindings.KeyMap) bool {
	return key.Matches(msg, keys.FocusDatabase) ||
		key.Matches(msg, keys.FocusSchemas) ||
		key.Matches(msg, keys.FocusQuery) ||
		key.Matches(msg, keys.FocusResults)
}

// HandleResultsTableKey handles keys when the Results panel is focused. Returns whether the key was consumed.
func (rt *ResultsTable) HandleResultsTableKey(msg tea.KeyMsg, keys bindings.KeyMap) (bool, tea.Cmd) {
	if rt == nil || !rt.active {
		return false, nil
	}

	if rt.editing {
		switch {
		case msg.Type == tea.KeyEnter:
			rt.commitEdit()
			return true, nil
		case key.Matches(msg, keys.Close) || msg.Type == tea.KeyEsc:
			rt.cancelEdit()
			return true, nil
		default:
			var c tea.Cmd
			rt.editor, c = rt.editor.Update(msg)
			return true, c
		}
	}

	if resultsKeyPassThrough(msg, keys) {
		return false, nil
	}

	nc := len(rt.colKeys)
	if nc == 0 {
		return false, nil
	}

	switch msg.Type {
	case tea.KeyTab:
		rt.col = (rt.col + 1) % nc
		return true, nil
	case tea.KeyShiftTab:
		rt.col = (rt.col - 1 + nc) % nc
		return true, nil
	case tea.KeyLeft:
		rt.col = (rt.col - 1 + nc) % nc
		return true, nil
	case tea.KeyRight:
		rt.col = (rt.col + 1) % nc
		return true, nil
	case tea.KeyEnter:
		return true, rt.startCellEdit()
	}

	if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == 'n' {
		rt.appendEmptyRow()
		return true, nil
	}

	// Wide table: horizontal scroll — handle types + [ ] explicitly (shift+arrows are unreliable in many terminals).
	if rt.kind == resultsKindWide {
		rt.tbl = rt.tbl.Focused(true)
		switch msg.Type {
		case tea.KeyShiftLeft:
			rt.tbl = rt.tbl.ScrollLeft()
			return true, nil
		case tea.KeyShiftRight:
			rt.tbl = rt.tbl.ScrollRight()
			return true, nil
		}
		if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
			switch msg.Runes[0] {
			case '[':
				rt.tbl = rt.tbl.ScrollLeft()
				return true, nil
			case ']':
				rt.tbl = rt.tbl.ScrollRight()
				return true, nil
			}
		}
	}

	rt.tbl = rt.tbl.Focused(true)
	var c tea.Cmd
	rt.tbl, c = rt.tbl.Update(msg)
	return true, c
}

func (rt *ResultsTable) startCellEdit() tea.Cmd {
	hi := rt.highlightedDataIndex()
	if hi < 0 || hi >= len(rt.data) {
		return nil
	}
	for len(rt.data[hi]) < len(rt.colKeys) {
		rt.data[hi] = append(rt.data[hi], "")
	}
	val := rt.data[hi][rt.col]
	rt.editor.SetValue(val)
	rt.editor.CursorEnd()
	rt.editing = true
	return rt.editor.Focus()
}

func (rt *ResultsTable) commitEdit() {
	if !rt.editing {
		return
	}
	hi := rt.highlightedDataIndex()
	if hi >= 0 && hi < len(rt.data) && rt.col < len(rt.colKeys) {
		for len(rt.data[hi]) < len(rt.colKeys) {
			rt.data[hi] = append(rt.data[hi], "")
		}
		rt.data[hi][rt.col] = rt.editor.Value()
	}
	rt.editing = false
	rt.editor.Blur()
}

func (rt *ResultsTable) cancelEdit() {
	if !rt.editing {
		return
	}
	rt.editing = false
	rt.editor.Blur()
}

// highlightedDataIndex maps the table's visible row cursor to our data slice (no filter / sort).
func (rt *ResultsTable) highlightedDataIndex() int {
	if len(rt.data) == 0 {
		return -1
	}
	hi := rt.tbl.GetHighlightedRowIndex()
	if hi < 0 {
		return 0
	}
	if hi >= len(rt.data) {
		return len(rt.data) - 1
	}
	return hi
}

func (rt *ResultsTable) appendEmptyRow() {
	nc := len(rt.colKeys)
	if nc == 0 {
		return
	}
	rt.data = append(rt.data, make([]string, nc))
	newHi := len(rt.data) - 1
	rt.tbl = rt.tbl.WithRows(rt.buildRowsForHi(newHi)).WithHighlightedRow(newHi)
	rt.clampCol()
}

func (rt *ResultsTable) clampCol() {
	if len(rt.colKeys) == 0 {
		rt.col = 0
		return
	}
	if rt.col < 0 {
		rt.col = 0
	}
	if rt.col >= len(rt.colKeys) {
		rt.col = len(rt.colKeys) - 1
	}
}

func (rt *ResultsTable) buildRowsForHi(highlightRow int) []table.Row {
	rows := make([]table.Row, 0, len(rt.data))
	for i, cells := range rt.data {
		rd := table.RowData{}
		for j, k := range rt.colKeys {
			val := ""
			if j < len(cells) {
				val = cells[j]
			}
			if i == highlightRow && j == rt.col && !rt.editing {
				rd[k] = table.NewStyledCell(val, cellFocusStyle)
			} else {
				rd[k] = val
			}
		}
		rd[colPadKey] = ""
		rows = append(rows, table.NewRow(rd))
	}
	return rows
}

func (rt *ResultsTable) syncRowsToTable(panelFocused bool) {
	if len(rt.data) == 0 {
		rt.tbl = rt.tbl.WithRows(nil)
		return
	}
	hi := rt.tbl.GetHighlightedRowIndex()
	if hi >= len(rt.data) {
		hi = len(rt.data) - 1
	}
	if hi < 0 {
		hi = 0
	}
	rows := rt.buildRowsForHi(hi)
	rt.tbl = rt.tbl.WithRows(rows).WithHighlightedRow(hi).Focused(panelFocused)
}

// RenderResultsTableBody returns inner content for draw.Border (hint + table [+ editor]).
func RenderResultsTableBody(rt *ResultsTable, innerW, innerH int, panelFocused bool) string {
	if rt == nil || !rt.active {
		return ""
	}
	if innerW < 8 {
		innerW = 8
	}
	if innerH < 1 {
		innerH = 1
	}

	rt.clampCol()
	rt.applyDimensions(innerW, innerH)
	rt.syncRowsToTable(panelFocused)

	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Faint(true).Width(innerW).Render(rt.hintLine()))
	b.WriteByte('\n')
	b.WriteString(rt.tbl.View())
	if rt.editing {
		b.WriteByte('\n')
		rt.editor.Width = innerW - 4
		if rt.editor.Width < 10 {
			rt.editor.Width = 10
		}
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("edit: "))
		b.WriteString(rt.editor.View())
	}
	return b.String()
}

func (rt *ResultsTable) hintLine() string {
	if rt == nil {
		return ""
	}
	if rt.editing {
		return "enter save · esc cancel"
	}
	base := "↑/↓ row · tab / shift+tab col · enter edit · n new row"
	switch rt.kind {
	case resultsKindTall:
		return base + " · h/l · pg page rows (footer)"
	case resultsKindWide:
		return base + " · [ / ] scroll columns (shift+←/→ if your terminal sends them)"
	default:
		return base
	}
}
