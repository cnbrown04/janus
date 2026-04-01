package ui

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cnbrown04/janus/bindings"
	"github.com/cnbrown04/janus/command"
	"github.com/cnbrown04/janus/draw"
)

// 1. Updated struct to hold the new nested layout fields
type Model struct {
	rootFlex  *flexbox.FlexBox
	rightFlex *flexbox.FlexBox
	rightCell *flexbox.Cell
	keys      bindings.KeyMap

	width, height int
	cmd            command.Model
	forceQuitArmed bool // first Ctrl+C arms; second Ctrl+C force-quits

	selectedPanel Panel

	queryArea           textarea.Model
	queryInsertMode     bool
	querySelAnchorLine  int // logical line index; -1 = no line selection
	queryScrollOff      int // first visible display row (soft-wrapped)
}

// New builds the UI model. A pointer is required so panel cells can read focus state while rendering.
func New() *Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.CharLimit = 0
	ta.Focus()

	m := &Model{
		keys:               bindings.DefaultKeyMap(),
		cmd:                command.New(),
		selectedPanel:      PanelQuery,
		queryArea:          ta,
		querySelAnchorLine: -1,
	}

	// Setup the Right Column (The nested FlexBox)
	right := flexbox.New(0, 0)

	// ratioY 1:4 gives Query a shorter strip; Results gets the rest.
	queryCell := flexbox.NewCell(1, 1).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		innerW := w - 2
		innerH := h - 2
		if innerW < 1 {
			innerW = 1
		}
		if innerH < 1 {
			innerH = 1
		}
		body := renderQueryBody(m, innerW, innerH)
		return draw.Border("Query", body, w, h, draw.PanelBorder(m.selectedPanel == PanelQuery))
	})
	dataCell := flexbox.NewCell(1, 4).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		return draw.Border("Results", "", w, h, draw.PanelBorder(m.selectedPanel == PanelResults))
	})
	
	right.AddRows([]*flexbox.Row{
		right.NewRow().AddCells(queryCell),
		right.NewRow().AddCells(dataCell),
	})

	// Setup the Root Layout
	root := flexbox.New(0, 0)
	
	leftSidebarCell := flexbox.NewCell(1, 1).SetStyle(lipgloss.NewStyle()).SetContentGenerator(func(w, h int) string {
		return draw.Border("Schemas", "", w, h, draw.PanelBorder(m.selectedPanel == PanelSchemas))
	})
	
	// This empty style is why lipgloss needs to be imported in this file
	rightContainerCell := flexbox.NewCell(6, 1).SetStyle(lipgloss.NewStyle())

	root.AddRows([]*flexbox.Row{
		root.NewRow().AddCells(leftSidebarCell, rightContainerCell),
	})

	m.rootFlex = root
	m.rightFlex = right
	m.rightCell = rightContainerCell
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// leaveQueryInsertMode exits INSERT without blurring the query textarea (normal mode still navigates).
func (m *Model) leaveQueryInsertMode() {
	m.queryInsertMode = false
	m.querySelAnchorLine = -1
}

// blurQueryPanel clears insert/selection state and blurs the query editor (other panels or modals).
func (m *Model) blurQueryPanel() {
	m.leaveQueryInsertMode()
	m.queryScrollOff = 0
	m.queryArea.Blur()
}
