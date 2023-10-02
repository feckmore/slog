package texthandler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// variation of https://github.com/jba/slog/blob/main/handlers/loghandler/log_handler.go

type Handler struct {
	opts      slog.HandlerOptions
	prefix    string // preformatted group names followed by a dot
	preformat string // preformatted Attrs, with an initial space

	mu sync.Mutex
	w  io.Writer
}

func New(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	h := &Handler{w: w}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.ReplaceAttr == nil {
		h.opts.ReplaceAttr = func(_ []string, a slog.Attr) slog.Attr { return a }
	}
	return h
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		w:         h.w,
		opts:      h.opts,
		preformat: h.preformat,
		prefix:    h.prefix + name + ".",
	}
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var buf []byte
	for _, a := range attrs {
		buf = h.appendAttr(buf, h.prefix, a)
	}
	return &Handler{
		w:         h.w,
		opts:      h.opts,
		prefix:    h.prefix,
		preformat: h.preformat + string(buf),
	}
}

const (
	whiteCircle  = "⚪"
	greenCircle  = "🟢"
	blueCircle   = "🔵"
	purpleCircle = "🟣"
	redCircle    = "🔴"
	orangeCircle = "🟠"
	yellowCircle = "🟡"
)

func appendLevel(buf []byte, level slog.Level) []byte {
	circles := map[slog.Level]string{
		slog.LevelDebug: whiteCircle,
		slog.LevelInfo:  blueCircle,
		slog.LevelWarn:  yellowCircle,
		slog.LevelError: redCircle,
	}

	circle := circles[level]
	if circle == "" {
		circle = purpleCircle
	}
	buf = append(buf, circle...)

	return fmt.Appendf(buf, " %-8s", level.String())
}

func appendFile(buf []byte, r slog.Record) []byte {
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		path := f.File
		filenameStart := strings.Index(path, "/slog/") + 6
		filename := path[filenameStart:]
		buf = append(buf, filename...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(f.Line), 10)
	}
	return buf
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var buf []byte
	buf = appendLevel(buf, r.Level)
	if !r.Time.IsZero() {
		buf = r.Time.AppendFormat(buf, time.TimeOnly)
		buf = append(buf, " | "...)
	}
	buf = append(buf, r.Message...)
	buf = append(buf, " | "...)
	buf = append(buf, h.preformat...)
	buf = append(buf, " | "...)
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, h.prefix, a)
		buf = append(buf, " | "...)
		return true
	})
	buf = appendFile(buf, r)
	buf = append(buf, '\n')
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf)
	return err
}

func (h *Handler) appendAttr(buf []byte, prefix string, a slog.Attr) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}
	if a.Value.Kind() != slog.KindGroup {
		buf = append(buf, ' ')
		buf = append(buf, prefix...)
		buf = append(buf, a.Key...)
		buf = append(buf, '=')
		return fmt.Appendf(buf, "%v", a.Value.Any())
	}
	// Group
	if a.Key != "" {
		prefix += a.Key + "."
	}
	for _, a := range a.Value.Group() {
		buf = h.appendAttr(buf, prefix, a)
	}
	return buf
}
