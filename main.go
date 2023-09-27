package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/feckmore/sandbox/slog/log"
)

const port = "8080"

func main() {
	log.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorldHandler)
	router := log.RequestIDMiddleWare(log.LoggingMiddleware(mux))

	slog.Info("listening", "port", port)
	err := http.ListenAndServe(":"+port, router)
	slog.Error("exiting", "error", err)
	os.Exit(1)
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Hello World Handler")
	fmt.Fprintf(w, "Hello World")
}
