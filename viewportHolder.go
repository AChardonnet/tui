package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewPortHolder struct {
	model viewport.Model
}

func (v *viewPortHolder) Init() tea.Cmd {
	return nil
}

func (v *viewPortHolder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	size, ok := msg.(tea.WindowSizeMsg)
	if ok {
		v.model.Width = size.Width
		v.model.Height = size.Height
		return v, nil
	}
	m, cmd := v.model.Update(msg)
	v.model = m
	return v, cmd
}

func (v *viewPortHolder) View() string {
	return v.model.View()
}
