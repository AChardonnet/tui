package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
	"golang.org/x/term"
)

var (
	lastWindowWidth  = -1
	lastWindowHeight = -1
)

type model struct {
	tui     boxer.Boxer
	logChan chan logMsg
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		waitForLog(m.logChan),
		windowSizePoll(),
	)
}

func windowSizePoll() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(250 * time.Millisecond)
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return nil
		}
		current := tea.WindowSizeMsg{Width: w, Height: h}
		lastWindowWidth = current.Width
		lastWindowHeight = current.Height
		return current
	}
}

func batchCmds(cmds ...tea.Cmd) tea.Cmd {
	filtered := make([]tea.Cmd, 0, len(cmds))
	for _, cmd := range cmds {
		if cmd != nil {
			filtered = append(filtered, cmd)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return tea.Batch(filtered...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	nextCmd := cmdForResize(msg)

	if size, ok := msg.(tea.WindowSizeMsg); ok {
		for _, leaf := range []string{mainAddr, playlistAddr, logsAddr} {
			if err := m.tui.EditLeaf(leaf, func(model tea.Model) (tea.Model, error) {
				updatedModel, _ := model.Update(msg)
				return updatedModel, nil
			}); err != nil {
				fmt.Printf("forward resize to leaf %s: %v\n", leaf, err)
			}
		}
		if err := m.tui.UpdateSize(size); err != nil {
			fmt.Printf("resize layout: %v\n", err)
		}
	}

	updatedTUI, cmd := m.tui.Update(msg)
	if updatedTUI == nil {
		updatedTUI = m.tui
	}

	tui, ok := updatedTUI.(boxer.Boxer)
	if !ok {
		return m, batchCmds(cmd, nextCmd)
	}

	m.tui = tui

	if logMsg, ok := msg.(logMsg); ok {
		if err := m.tui.EditLeaf(logsAddr, func(model tea.Model) (tea.Model, error) {
			lm, ok := model.(*logModel)
			if !ok {
				return model, fmt.Errorf("log leaf is %T", model)
			}
			updated, _ := lm.Update(logMsg)
			return updated, nil
		}); err != nil {
			fmt.Printf("edit log leaf: %v\n", err)
		}
		return m, batchCmds(cmd, waitForLog(m.logChan), nextCmd)
	}

	return m, batchCmds(cmd, nextCmd)
}

func cmdForResize(msg tea.Msg) tea.Cmd {
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		return windowSizePoll()
	}
	return nil
}

func (m model) View() string {
	return m.tui.View()
}

func (m *model) editModel(addr string, edit func(tea.Model) (tea.Model, error)) error {
	if edit == nil {
		return fmt.Errorf("no edit function provided")
	}
	v, ok := m.tui.ModelMap[addr]
	if !ok {
		return fmt.Errorf("no model with address '%s' found", addr)
	}
	v, err := edit(v)
	if err != nil {
		return err
	}
	m.tui.ModelMap[addr] = v
	return nil
}
