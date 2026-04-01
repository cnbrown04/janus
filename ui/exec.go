package ui

import "github.com/cnbrown04/janus/ui/modules"

func (m *Model) runQueryExecute() {
	text := modules.QueryExecuteText(m.queryArea.Value(), m.querySelAnchorLine, m.queryArea.Line())
	_ = text // TODO: execute against connection / results panel
	m.querySelAnchorLine = -1
}
