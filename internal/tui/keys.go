package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding // collapse dir
	Right      key.Binding // expand dir
	PageUp     key.Binding
	PageDown   key.Binding
	Toggle     key.Binding
	Expand     key.Binding
	SelectAll  key.Binding
	SelectAllR key.Binding
	Sort       key.Binding
	ToggleGit  key.Binding
	ToggleHide key.Binding
	ToggleHex  key.Binding
	ToggleExt  key.Binding
	ResetExt   key.Binding
	Tab        key.Binding
	Confirm    key.Binding
	Quit       key.Binding
	Help       key.Binding
}

var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:       key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "collapse")),
	Right:      key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "expand")),
	PageUp:     key.NewBinding(key.WithKeys("pgup", "ctrl+u"), key.WithHelp("pgup/^U", "page up")),
	PageDown:   key.NewBinding(key.WithKeys("pgdown", "ctrl+d"), key.WithHelp("pgdn/^D", "page dn")),
	Toggle:     key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
	Expand:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "collapse/expand")),
	SelectAll:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select dir")),
	SelectAllR: key.NewBinding(key.WithKeys("A"), key.WithHelp("A", "select all")),
	Sort:       key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "sort")),
	ToggleGit:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "gitignore")),
	ToggleHide: key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "hidden files")),
	ToggleHex:  key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "hex dump")),
	ToggleExt:  key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "toggle ext")),
	ResetExt:   key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "reset ext")),
	Tab:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
	Confirm:    key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("ctrl+g", "generate")),
	Quit:       key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "quit")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}
