package ui

// Panel identifies the main content panes (Database, Schemas, Query, Results).
type Panel uint8

const (
	PanelDatabase Panel = iota
	PanelSchemas
	PanelQuery
	PanelResults
)
