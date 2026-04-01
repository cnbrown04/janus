package ui

// Panel identifies one of the three main content panes.
type Panel uint8

const (
	PanelSchemas Panel = iota
	PanelQuery
	PanelResults
)
