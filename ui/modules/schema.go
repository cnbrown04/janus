package modules

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/draw"
)

// testCatalog is stub schema/table data for the Schemas panel until a real connection exists.
var testCatalog = struct {
	schema string
	tables []string
}{
	schema: "test",
	tables: []string{"user", "tall", "wide"},
}

// SchemaTreeState holds expand/collapse and cursor for the schema sidebar (neo-tree–style navigation).
type SchemaTreeState struct {
	Expanded bool
	Cursor   int
}

// SchemaTreeResult is the outcome of handling a key in the schema tree.
type SchemaTreeResult struct {
	Handled   bool
	OpenTable string // non-empty → open this table in the Results panel
}

func catalogLineCount(expanded bool) int {
	if !expanded {
		return 1
	}
	return 1 + len(testCatalog.tables)
}

func (st *SchemaTreeState) clampCursor() {
	max := catalogLineCount(st.Expanded) - 1
	if max < 0 {
		max = 0
	}
	if st.Cursor < 0 {
		st.Cursor = 0
	}
	if st.Cursor > max {
		st.Cursor = max
	}
}

// SchemaCatalogName returns the stub schema name (for previews and tests).
func SchemaCatalogName() string {
	return testCatalog.schema
}

// HandleSchemaTreeKey handles neo-tree–like navigation: ↑/k, ↓/j, Enter (toggle folder or open table).
func HandleSchemaTreeKey(msg tea.KeyMsg, keys bindings.KeyMap, st *SchemaTreeState) SchemaTreeResult {
	if key.Matches(msg, keys.Up) {
		st.Cursor--
		st.clampCursor()
		return SchemaTreeResult{Handled: true}
	}
	if key.Matches(msg, keys.Down) {
		st.Cursor++
		st.clampCursor()
		return SchemaTreeResult{Handled: true}
	}
	if msg.Type == tea.KeyEnter {
		if st.Cursor == 0 {
			st.Expanded = !st.Expanded
			if !st.Expanded {
				st.Cursor = 0
			}
			st.clampCursor()
			return SchemaTreeResult{Handled: true}
		}
		if st.Expanded && st.Cursor >= 1 {
			idx := st.Cursor - 1
			if idx < len(testCatalog.tables) {
				return SchemaTreeResult{Handled: true, OpenTable: testCatalog.tables[idx]}
			}
		}
		return SchemaTreeResult{Handled: true}
	}
	return SchemaTreeResult{}
}

func schemaScrollOffset(cursor, viewH, total int) int {
	if total <= viewH {
		return 0
	}
	maxScroll := total - viewH
	s := cursor - viewH + 1
	if s < 0 {
		s = 0
	}
	if s > maxScroll {
		s = maxScroll
	}
	return s
}

var schemaSelStyle = lipgloss.NewStyle().
	Background(draw.SelectionBG).
	Foreground(draw.SelectionFG)

func formatTreeLine(plain string, selected bool, innerW int) string {
	padded := PadLineToWidth(plain, innerW)
	if selected {
		return schemaSelStyle.Render(padded)
	}
	return padded
}

// schemaTreeLines builds visible tree lines (schema row + optional table rows).
func schemaTreeLines(innerW int, focused bool, st SchemaTreeState) []string {
	discl := "▸ "
	if st.Expanded {
		discl = "▼ "
	}
	schemaLine := discl + testCatalog.schema

	var lines []string
	sel0 := focused && st.Cursor == 0
	lines = append(lines, formatTreeLine(schemaLine, sel0, innerW))

	if !st.Expanded {
		return lines
	}
	for i, tbl := range testCatalog.tables {
		prefix := "├── "
		if i == len(testCatalog.tables)-1 {
			prefix = "└── "
		}
		sel := focused && st.Cursor == 1+i
		lines = append(lines, formatTreeLine(prefix+tbl, sel, innerW))
	}
	return lines
}

// RenderSchemasPanel draws the Schemas sidebar with a Unicode tree and selection highlight.
func RenderSchemasPanel(w, h int, focused bool, st SchemaTreeState) string {
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	innerLines := h - 2
	if innerLines < 1 {
		innerLines = 1
	}

	all := schemaTreeLines(innerW, focused, st)
	total := len(all)
	scroll := schemaScrollOffset(st.Cursor, innerLines, total)
	end := scroll + innerLines
	if end > len(all) {
		end = len(all)
	}
	window := all[scroll:end]
	for len(window) < innerLines {
		window = append(window, strings.Repeat(" ", innerW))
	}
	body := strings.Join(window, "\n")
	return draw.Border("Schemas", body, w, h, draw.PanelBorder(focused))
}
