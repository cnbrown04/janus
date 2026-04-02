package ui

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/command"
	"github.com/cnbrown04/janus/draw"
	"github.com/cnbrown04/janus/ui/modules"
)

// Model holds the full UI layout and editor state.
type Model struct {
	rootFlex  *flexbox.FlexBox
	leftFlex  *flexbox.FlexBox
	rightFlex *flexbox.FlexBox
	rightCell *flexbox.Cell
	keys      bindings.KeyMap

	width, height  int
	cmd            command.Model
	forceQuitArmed bool

	selectedPanel Panel

	queryArea          textarea.Model
	queryInsertMode    bool
	querySelAnchorLine int
	queryScrollOff     int

	Database     modules.DatabasePanelState
	SchemaTree   modules.SchemaTreeState
	ResultsBody  string
	ResultsTable modules.ResultsTable
}

// New builds the UI model. A pointer is required so panel cells can read focus state while rendering.
func New() *Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.CharLimit = 0
	ta.Focus()

	dbOpts := modules.DefaultDatabaseOptions()
	m := &Model{
		keys:               bindings.DefaultKeyMap(),
		cmd:                command.New(),
		selectedPanel:      PanelQuery,
		queryArea:          ta,
		querySelAnchorLine: -1,
		Database: modules.DatabasePanelState{
			Options:  dbOpts,
			Selected: dbOpts[0],
		},
	}

	right := flexbox.New(0, 0)

	queryCell := flexbox.NewCell(1, 1).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		innerW := w - 2
		innerH := h - 2
		if innerW < 1 {
			innerW = 1
		}
		if innerH < 1 {
			innerH = 1
		}
		body := modules.RenderQueryBody(&modules.QueryRenderContext{
			TA:            &m.queryArea,
			SelAnchorLine: m.querySelAnchorLine,
			ScrollOff:     &m.queryScrollOff,
		}, innerW, innerH)
		return draw.Border("Query", body, w, h, draw.PanelBorder(m.selectedPanel == PanelQuery))
	})
	dataCell := flexbox.NewCell(1, 4).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		return modules.RenderResultsPanel(w, h, m.selectedPanel == PanelResults, m.ResultsBody, &m.ResultsTable)
	})

	right.AddRows([]*flexbox.Row{
		right.NewRow().AddCells(queryCell),
		right.NewRow().AddCells(dataCell),
	})

	root := flexbox.New(0, 0)

	left := flexbox.New(0, 0)
	topCell := flexbox.NewCell(1, 1).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		return modules.RenderDatabasePanel(w, h, &m.Database, m.selectedPanel == PanelDatabase)
	})
	schemasCell := flexbox.NewCell(1, 17).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		return modules.RenderSchemasPanel(w, h, m.selectedPanel == PanelSchemas, m.SchemaTree)
	})
	left.AddRows([]*flexbox.Row{
		left.NewRow().AddCells(topCell),
		left.NewRow().AddCells(schemasCell),
	})

	leftSidebarCell := flexbox.NewCell(1, 1).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		m.leftFlex.SetWidth(w)
		m.leftFlex.SetHeight(h)
		return m.leftFlex.Render()
	})

	rightContainerCell := flexbox.NewCell(6, 1).SetStyle(lipgloss.NewStyle())

	root.AddRows([]*flexbox.Row{
		root.NewRow().AddCells(leftSidebarCell, rightContainerCell),
	})

	m.rootFlex = root
	m.leftFlex = left
	m.rightFlex = right
	m.rightCell = rightContainerCell
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) leaveQueryInsertMode() {
	m.queryInsertMode = false
	m.querySelAnchorLine = -1
}

func (m *Model) blurQueryPanel() {
	m.leaveQueryInsertMode()
	m.queryScrollOff = 0
	m.queryArea.Blur()
	m.Database.DropdownOpen = false
	m.ResultsTable.Blur()
}

func (m *Model) handleSchemaTreeKey(msg tea.KeyMsg) bool {
	r := modules.HandleSchemaTreeKey(msg, m.keys, &m.SchemaTree)
	if r.OpenTable != "" {
		m.blurQueryPanel()
		if rt := modules.NewResultsTableForCatalog(r.OpenTable); rt.Active() {
			m.ResultsTable = rt
			m.ResultsBody = ""
		} else {
			m.ResultsTable = modules.ResultsTable{}
			m.ResultsBody = modules.FormatTablePreview(modules.SchemaCatalogName(), r.OpenTable)
		}
		m.selectedPanel = PanelResults
		return true
	}
	return r.Handled
}
