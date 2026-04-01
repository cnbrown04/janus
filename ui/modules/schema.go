package modules

import "github.com/cnbrown04/janus/draw"

// RenderSchemasPanel draws the Schemas sidebar panel.
func RenderSchemasPanel(w, h int, focused bool) string {
	return draw.Border("Schemas", "", w, h, draw.PanelBorder(focused))
}
