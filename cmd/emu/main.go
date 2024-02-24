package main

import (
	"gogb/pkg/cpu"
	"gogb/pkg/mem"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func main() {

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))

	slog.Info("Hello!")
	if len(os.Args) < 2 {
		slog.Info("usage: emu [bin-path]")
		os.Exit(1)
	}
	path := os.Args[1]
	bytes, err := os.ReadFile(path)
	if err != nil {
		slog.Error("err: %s", err)
		os.Exit(1)
	}

	mem := mem.NewRAM()
	mem.Copy(bytes, 0)

	cpu := cpu.NewCPU(&mem)

	for {
		cpu.FetchExecute()
		// time.Sleep(time.Second / 100)
	}
}
