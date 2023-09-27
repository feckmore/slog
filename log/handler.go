package log

import (
	"context"
	"log/slog"
)

type LogHandler struct {
	slog.Handler
}

// Handle inserts the request_id parameter into each log record if found in the context
// and is required to implement slog.Handler
func (h *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	if id, exists := RequestID(ctx); exists {
		record.AddAttrs(slog.String("request_id", id))
	}

	return h.Handler.Handle(ctx, record)
}

func (h *LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{Handler: h.Handler.WithGroup(name)}
}
