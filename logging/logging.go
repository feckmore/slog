package logging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/feckmore/sandbox/slog/requestid"
)

type LogHandler struct {
	slog.Handler
}

func Initialize() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	myHandler := LogHandler{
		Handler: jsonHandler,
	}
	logger := slog.New(&myHandler)
	slog.SetDefault(logger)
	slog.Info("logger initialized")
}

// Middleware logs the start and end of the request, and includes the time to process
func Middleware(next http.Handler) http.Handler {
	anonymousFunction := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.DebugContext(r.Context(), "request start")
		next.ServeHTTP(w, r)
		stop := time.Now()
		slog.DebugContext(
			r.Context(),
			"request end",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.String("duration", fmt.Sprintf("%d μs", stop.Sub(start).Microseconds())),
		)
	}

	return http.HandlerFunc(anonymousFunction)
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, exists := requestid.RequestID(ctx); exists {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
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
