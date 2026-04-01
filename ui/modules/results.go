package modules

import (
	"strings"

	"github.com/cnbrown04/janus/draw"
)

// RenderResultsPanel draws the Results panel.
func RenderResultsPanel(w, h int, focused bool, body string) string {
	body = strings.TrimRight(body, "\n")
	return draw.Border("Results", body, w, h, draw.PanelBorder(focused))
}
