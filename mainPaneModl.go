package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	viewport viewport.Model
	cacheBar CacheBar
}

func newMainModel() *mainModel {
	return &mainModel{
		viewport: viewport.New(0, 0),
		cacheBar: *newCacheBar(),
	}
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) View() string {
	content := m.cacheBar.View() + "\n\n"
	content += "Some other content here\n"
	content += "More content\n"
	return content
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.cacheBar.Width = size.Width
	}
	return m, nil
}
