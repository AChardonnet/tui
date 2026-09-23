package main

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BlockState int

const (
	BlockEmpty BlockState = iota
	BlockLoading
	BlockCached
	BlockFailed
)

var (
	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220"))

	cachedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	failedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type CacheBlock struct {
	Start float64
	End   float64
	State BlockState
}

type CacheBar struct {
	Width  int
	Blocks []CacheBlock
}

func newCacheBar() *CacheBar {
	var blocks []CacheBlock

	for i := 0; i < 100; i++ {
		blocks = append(blocks, CacheBlock{
			Start: float64(i * 1000),
			End:   float64((i + 1) * 1000),
			State: BlockState(rand.Intn(4)),
		})
	}
	return &CacheBar{Width: 80, Blocks: blocks}
}

func (b CacheBar) Init() tea.Cmd {
	return nil
}

func (b CacheBar) View() string {
	if b.Width <= 0 {
		return ""
	}

	if len(b.Blocks) == 0 {
		return strings.Repeat(emptyStyle.Render("░"), b.Width)
	}

	var s strings.Builder
	for i := 0; i < b.Width; i++ {
		idx := int(float64(i) * float64(len(b.Blocks)) / float64(b.Width))
		if idx >= len(b.Blocks) {
			idx = len(b.Blocks) - 1
		}

		switch b.Blocks[idx].State {
		case BlockCached:
			s.WriteString(cachedStyle.Render("█"))
		case BlockLoading:
			s.WriteString(loadingStyle.Render("█"))
		case BlockFailed:
			s.WriteString(failedStyle.Render("█"))
		default:
			s.WriteString(emptyStyle.Render("░"))
		}
	}

	return s.String()
}

func (b CacheBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.Width = msg.Width
		return b, nil
	default:
		return b, nil
	}
}
