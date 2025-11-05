package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit       key.Binding
	Refresh    key.Binding
	Stage      key.Binding
	Unstage    key.Binding
	Commit     key.Binding
	Push       key.Binding
	Pull       key.Binding
	Up         key.Binding
	Down       key.Binding
	ToggleDiff key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Refresh:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Stage:      key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "stage")),
		Unstage:    key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "unstage")),
		Commit:     key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "commit")),
		Push:       key.NewBinding(key.WithKeys("P"), key.WithHelp("P", "push")),
		Pull:       key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "pull/fetch")),
		Up:         key.NewBinding(key.WithKeys("k", "up")),
		Down:       key.NewBinding(key.WithKeys("j", "down")),
		ToggleDiff: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle diff")),
	}
}
