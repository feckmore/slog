package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/feckmore/sandbox/slog/logging"
	"github.com/feckmore/sandbox/slog/requestid"
)

const port = "8080"

func handler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Hello World Handler")
	fmt.Fprintf(w, "Hello World")
}

func main() {
	logging.Initialize()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	router := requestid.MiddleWare(logging.Middleware(mux))

	slog.Debug("listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
