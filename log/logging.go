package log

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Init() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	myHandler := LogHandler{
		Handler: jsonHandler,
	}
	logger := slog.New(&myHandler)
	slog.SetDefault(logger)
	slog.Info("logger initialized")
}

// LoggingMiddleware logs the start and end of the request, and includes the time to process
func LoggingMiddleware(next http.Handler) http.Handler {
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
