package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type logLevel int

const (
	levelDebug logLevel = iota
	levelInfo
	levelWarn
	levelError
)

var loggingLevel logLevel = levelDebug

const (
	maxLogs = 1000
)

type logMsg struct {
	Level   slog.Level
	Message string
	Time    time.Time
}

type logItem struct {
	Level   slog.Level
	Message string
	Time    time.Time
	Line    string
}

func newLogItemFromLogMessage(logMsg logMsg) logItem {
	line := fmt.Sprintf(
		"[%s] [%s] %s",
		logMsg.Time.Format("15:04:05"),
		logMsg.Level.String(),
		logMsg.Message,
	)
	return logItem{
		Level:   logMsg.Level,
		Message: logMsg.Message,
		Time:    logMsg.Time,
		Line:    line,
	}
}

func isLogValid(level slog.Level) bool {
	var logLvl logLevel
	switch level {
	case slog.LevelDebug:
		logLvl = levelDebug
	case slog.LevelInfo:
		logLvl = levelInfo
	case slog.LevelWarn:
		logLvl = levelWarn
	case slog.LevelError:
		logLvl = levelError
	default:
		logLvl = levelDebug
	}
	return logLvl >= loggingLevel
}

func setupLogger(ch chan<- logMsg) (*os.File, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(
		"logs/log.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	fileHandler := slog.NewTextHandler(
		logFile,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)

	tuiHandler := &tuiLogHandler{
		ch: ch,
	}

	logger := slog.New(&multiHandler{
		handlers: []slog.Handler{
			fileHandler,
			tuiHandler,
		},
	})

	slog.SetDefault(logger)

	return logFile, nil
}

func waitForLog(ch <-chan logMsg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
