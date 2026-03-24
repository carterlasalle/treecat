package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Toggle     key.Binding
	Expand     key.Binding
	SelectAll  key.Binding
	SelectAllR key.Binding
	Sort       key.Binding
	ToggleGit  key.Binding
	ToggleHide key.Binding
	ToggleHex  key.Binding
	Tab        key.Binding
	Confirm    key.Binding
	Quit       key.Binding
	Help       key.Binding
}

var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Toggle:     key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
	Expand:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "expand/collapse")),
	SelectAll:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all (dir)")),
	SelectAllR: key.NewBinding(key.WithKeys("A"), key.WithHelp("A", "select all (recursive)")),
	Sort:       key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "cycle sort")),
	ToggleGit:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "toggle gitignore")),
	ToggleHide: key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "toggle hidden")),
	ToggleHex:  key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "hex binary")),
	Tab:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
	// ctrl+c is intercepted by Bubble Tea before reaching the model; use ctrl+g to confirm
	Confirm: key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("ctrl+g", "generate output")),
	Quit:    key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "quit")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}
