package log

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/feckmore/sandbox/slog/log/loghandler"
	"github.com/feckmore/sandbox/slog/log/texthandler"
)

func Init() {
	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = texthandler.New(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = loghandler.New(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(handler)
	logger = logger.With(slog.String("app", "slog"), slog.String("version", "1.0.0"), slog.String("env", "dev"))
	slog.SetDefault(logger)

	slog.Debug("example debug log")
	slog.Info("example info log")
	slog.Warn("example warn log")
	slog.Error("example error log")
}

// LoggingMiddleware logs the start and end of the request, and includes the time to process
func LoggingMiddleware(next http.Handler) http.Handler {
	anonymousFunction := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.DebugContext(
			r.Context(),
			"request start",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)
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
