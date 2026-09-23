package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestCacheBarWidth(t *testing.T) {
	b := newCacheBar()

	updated, _ := b.Update(tea.WindowSizeMsg{Width: 69, Height: 42})
	updatedBar, ok := updated.(CacheBar)
	if !ok {
		t.Fatalf("cache bar update returned %T, want CacheBar", updated)
	}

	if updatedBar.Width != 69 {
		t.Fatalf("bar width did not update correctly")
	}
}

func TestMainModelResizeKeepsViewWithinPaneWidth(t *testing.T) {
	m := newMainModel()

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 76, Height: 20})
	updatedModel, ok := updated.(*mainModel)
	if !ok {
		t.Fatalf("main model update returned %T, want *mainModel", updated)
	}

	for _, line := range strings.Split(updatedModel.View(), "\n") {
		if width := lipgloss.Width(line); width > 76 {
			t.Fatalf("main view line width is %d, want at most 76", width)
		}
	}
}
