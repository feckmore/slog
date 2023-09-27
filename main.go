package main

import (
	"fmt"
	l "log"
	"log/slog"
	"net/http"

	"github.com/feckmore/sandbox/slog/log"
)

const port = "8080"

func handler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Hello World Handler")
	fmt.Fprintf(w, "Hello World")
}

func main() {
	log.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	router := log.RequestIDMiddleWare(log.LoggingMiddleware(mux))

	slog.Debug("listening on port " + port)

	l.Fatal(http.ListenAndServe(":"+port, router))
}
