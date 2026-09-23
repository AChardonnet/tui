package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type logModel struct {
	model       viewport.Model
	logs        []logItem
	printedLogs []string
}

func newLogModel() *logModel {
	return &logModel{
		model:       viewport.New(0, 0),
		logs:        []logItem{},
		printedLogs: []string{},
	}
}

func (m *logModel) Init() tea.Cmd {
	return nil
}

func (m *logModel) Add(msg logMsg) {
	logI := newLogItemFromLogMessage(msg)
	m.logs = append(m.logs, logI)

	if isLogValid(logI.Level) {
		m.printedLogs = append(m.printedLogs, logI.Line)
	}

	if len(m.logs) > maxLogs {
		m.logs = m.logs[len(m.logs)-maxLogs:]
	}

	if len(m.printedLogs) > maxLogs {
		m.printedLogs = m.printedLogs[len(m.printedLogs)-maxLogs:]
	}

	m.model.SetContent(wrapLogLines(m.printedLogs, m.model.Width))
	m.model.GotoBottom()
}

func wrapLogLines(lines []string, width int) string {
	if width <= 0 {
		return strings.Join(lines, "\n")
	}

	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			wrapped = append(wrapped, "")
			continue
		}
		wrapped = append(wrapped, lipgloss.NewStyle().Width(width).Render(line))
	}
	return strings.Join(wrapped, "\n")
}

func (m *logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.Width = msg.Width
		m.model.Height = msg.Height
		m.model.SetContent(wrapLogLines(m.printedLogs, m.model.Width))
		return m, nil

	case logMsg:
		m.Add(msg)
		return m, nil
	}

	newModel, cmd := m.model.Update(msg)
	m.model = newModel

	return m, cmd
}

func (m *logModel) View() string {
	return m.model.View()
}
