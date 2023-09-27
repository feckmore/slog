package log

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"log/slog"

	"github.com/oklog/ulid"
)

type contextKey string

const contextKeyRequestID contextKey = "request-id-key"

// RequestID extracts the request ID from a context if present
func RequestID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	id, exists := ctx.Value(contextKeyRequestID).(string)
	return id, exists
}

// RequestIDMiddleWare adds a unique ID to the request's context
func RequestIDMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(SetRequestID(r.Context(), NewRequestID())))
	})
}

// NewRequestID returns a new ULID (https://github.com/ulid/spec)
func NewRequestID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	requestID, err := ulid.New(ms, entropy)
	if err != nil {
		slog.Error("error creating request id", "error", err)
	}

	return requestID.String()
}

// SetRequestID includes the request ID in the request context
func SetRequestID(parent context.Context, requestID string) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	return context.WithValue(parent, contextKeyRequestID, requestID)
}
