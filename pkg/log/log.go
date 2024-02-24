package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var LogLevel = new(slog.LevelVar)

func Init() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: LogLevel})))
}
