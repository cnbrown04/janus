package bindings

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all the keybindings for the app.
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Command key.Binding // toggles command bar (e.g. ":")
	Close   key.Binding // closes command bar

	FocusDatabase key.Binding
	FocusSchemas    key.Binding
	FocusQuery      key.Binding
	FocusResults    key.Binding

	QueryInsert key.Binding // i — insert mode in Query panel

	QueryExecute key.Binding // ctrl+j / alt+enter — run query (selection or all)
}

// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Command: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "command"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close"),
		),
		FocusDatabase: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "database"),
		),
		FocusSchemas: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "schemas"),
		),
		FocusQuery: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "query"),
		),
		FocusResults: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "results"),
		),
		QueryInsert: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert (query)"),
		),
		QueryExecute: key.NewBinding(
			key.WithKeys("ctrl+j", "alt+enter"),
			key.WithHelp("ctrl+j", "run query"),
		),
	}
}
