package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/lmittmann/tint"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, nil)))

	http.Handle("/", http.FileServer(http.Dir("web/public")))
	slog.Info("Listening and serving @ :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
