package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cnbrown04/janus/command"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.rootFlex.SetWidth(msg.Width)
		m.rootFlex.SetHeight(msg.Height)

		rightSideWidth := (msg.Width * 6) / 7
		m.rightFlex.SetWidth(rightSideWidth)
		m.rightFlex.SetHeight(msg.Height)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			if m.forceQuitArmed {
				return m, tea.Quit
			}
			m.forceQuitArmed = true
			return m, nil
		}
		m.forceQuitArmed = false

		if m.cmd.Active() && key.Matches(msg, m.keys.Close) {
			m.cmd = m.cmd.Close()
			return m, nil
		}

		if key.Matches(msg, m.keys.Command) {
			m.blurQueryPanel()
			var c tea.Cmd
			m.cmd, c = m.cmd.Toggle()
			return m, c
		}

		if m.cmd.Active() {
			var c tea.Cmd
			m.cmd, c = command.Update(m.cmd, msg)
			return m, c
		}

		if m.selectedPanel == PanelQuery {
			if m.matchesQueryExecute(msg) {
				m.runQueryExecute()
				return m, nil
			}

			if m.queryInsertMode {
				if key.Matches(msg, m.keys.Close) {
					m.leaveQueryInsertMode()
					return m, nil
				}

				switch msg.Type {
				case tea.KeyShiftUp, tea.KeyShiftDown:
					if m.querySelAnchorLine < 0 {
						m.querySelAnchorLine = m.queryArea.Line()
					}
				case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
					m.querySelAnchorLine = -1
				default:
					if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
						m.querySelAnchorLine = -1
					}
				}

				fwd := remapShiftVerticalForTextarea(msg)

				var c tea.Cmd
				m.queryArea, c = m.queryArea.Update(fwd)
				return m, c
			}

			// Normal mode on Query: navigate and line-select without inserting text.
			if key.Matches(msg, m.keys.QueryInsert) {
				m.queryInsertMode = true
				return m, m.queryArea.Focus()
			}

			if shouldForwardQueryNormalMode(msg) {
				switch msg.Type {
				case tea.KeyShiftUp, tea.KeyShiftDown:
					if m.querySelAnchorLine < 0 {
						m.querySelAnchorLine = m.queryArea.Line()
					}
				case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
					m.querySelAnchorLine = -1
				default:
					if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
						m.querySelAnchorLine = -1
					}
				}

				fwd := remapShiftVerticalForTextarea(msg)

				var c tea.Cmd
				m.queryArea, c = m.queryArea.Update(fwd)
				return m, c
			}
		}

		if !m.queryInsertMode {
			switch {
			case key.Matches(msg, m.keys.FocusSchemas):
				m.blurQueryPanel()
				m.selectedPanel = PanelSchemas
			case key.Matches(msg, m.keys.FocusQuery):
				m.selectedPanel = PanelQuery
				return m, m.queryArea.Focus()
			case key.Matches(msg, m.keys.FocusResults):
				m.blurQueryPanel()
				m.selectedPanel = PanelResults
			}
		}
	}

	return m, nil
}
