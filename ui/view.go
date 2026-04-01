package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	// Reserve one line per footer so the main flex layout stays within m.height total.
	footerLines := 0
	if m.selectedPanel == PanelQuery && !m.cmd.Active() && m.width > 0 {
		footerLines++
	}
	if m.forceQuitArmed {
		footerLines++
	}

	contentH := m.height - footerLines
	if contentH < 1 {
		contentH = 1
	}
	if m.width > 0 && m.height > 0 {
		m.rootFlex.SetHeight(contentH)
		m.rightFlex.SetHeight(contentH)
	}

	rightColumnString := m.rightFlex.Render()
	m.rightCell.SetContent(rightColumnString)
	layout := m.rootFlex.Render()

	if m.cmd.Active() && m.width > 0 && contentH > 0 {
		if dlg := m.cmd.View(m.width); dlg != "" {
			layout = overlayModal(layout, dlg, m.width, contentH)
		}
	}

	if m.selectedPanel == PanelQuery && !m.cmd.Active() && m.width > 0 {
		var footer string
		if m.queryInsertMode {
			footer = lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.Color("244")).
				Width(m.width).
				Render("-- INSERT --   esc normal   ·   ctrl+j / alt+↵ run (selection or all)   ·   shift+↑/↓ line selection")
		} else {
			footer = lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.Color("244")).
				Width(m.width).
				Render("i insert   ·   shift+↑/↓ select lines   ·   ctrl+j / alt+↵ run (selection or all)   ·   ⌘B bind param (stub)")
		}
		layout = lipgloss.JoinVertical(lipgloss.Top, layout, footer)
	}

	if m.forceQuitArmed {
		hint := lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.Color("241")).
			Render("Press Ctrl+C again to force quit.")
		layout = lipgloss.JoinVertical(lipgloss.Top, layout, hint)
	}

	return layout
}
