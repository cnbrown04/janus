package modules

import "github.com/cnbrown04/janus/draw"

// RenderResultsPanel draws the Results panel.
func RenderResultsPanel(w, h int, focused bool) string {
	return draw.Border("Results", "", w, h, draw.PanelBorder(focused))
}
