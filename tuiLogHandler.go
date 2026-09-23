package main

import (
	"context"
	"log/slog"
)

type tuiLogHandler struct {
	ch chan<- logMsg
}

func (h *tuiLogHandler) Enabled(
	_ context.Context,
	level slog.Level,
) bool {
	return true
}

func (h *tuiLogHandler) Handle(
	_ context.Context,
	record slog.Record,
) error {
	h.ch <- logMsg{
		Level:   record.Level,
		Message: record.Message,
		Time:    record.Time,
	}

	return nil
}

func (h *tuiLogHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {
	return h
}

func (h *tuiLogHandler) WithGroup(
	name string,
) slog.Handler {
	return h
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(
	ctx context.Context,
	level slog.Level,
) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (h *multiHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *multiHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))

	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}

	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(
	name string,
) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))

	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}

	return &multiHandler{handlers: handlers}
}
