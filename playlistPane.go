package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type playlistItem struct {
	Title  string
	Status string
}

type playlistModel struct {
	items    []playlistItem
	selected int
	width    int
	height   int
}

func newPlaylistModel() *playlistModel {
	return &playlistModel{
		items: []playlistItem{
			{Title: "Intro", Status: "ready"},
			{Title: "Main topic", Status: "downloading"},
			{Title: "Outro", Status: "queued"},
		},
		selected: 1,
	}
}

func (m *playlistModel) Init() tea.Cmd {
	return nil
}

func (m *playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

func (m *playlistModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	placeholderStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("245"))

	var b strings.Builder
	b.WriteString(titleStyle.Render("Playlist"))
	b.WriteString("\n")

	if len(m.items) == 0 {
		b.WriteString(placeholderStyle.Render("No videos in queue"))
		return b.String()
	}

	for i, item := range m.items {
		line := fmt.Sprintf("%s %s", prefixForSelection(i == m.selected), item.Title)
		if i == m.selected {
			b.WriteString(selectedStyle.Render(line))
		} else {
			b.WriteString(line)
		}

		if item.Status != "" {
			b.WriteString(mutedStyle.Render("  [" + item.Status + "]"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func prefixForSelection(selected bool) string {
	if selected {
		return "▶"
	}
	return "•"
}
