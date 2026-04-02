package modules

import (
	"fmt"
	"strings"

	"github.com/cnbrown04/janus/draw"
)

// RenderResultsPanel draws the Results panel (plain body and/or interactive bubble-table).
func RenderResultsPanel(w, h int, focused bool, body string, rt *ResultsTable) string {
	innerW := w - 2
	if innerW < 1 {
		innerW = 1
	}
	innerH := h - 2
	if innerH < 1 {
		innerH = 1
	}
	if rt != nil && rt.Active() {
		content := RenderResultsTableBody(rt, innerW, innerH, focused)
		return draw.Border("Results", content, w, h, draw.PanelBorder(focused))
	}
	body = strings.TrimRight(body, "\n")
	return draw.Border("Results", body, w, h, draw.PanelBorder(focused))
}

// FormatTablePreview returns Results body text for catalog tables without an interactive bubble-table.
func FormatTablePreview(schema, table string) string {
	return fmt.Sprintf("Table %s.%s\n\n(Not connected — preview only.)", schema, table)
}
