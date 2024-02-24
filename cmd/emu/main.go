package main

import (
	"gogb/pkg/cpu"
	"gogb/pkg/log"
	"gogb/pkg/mem"
	"log/slog"
	"os"
)

func main() {

	log.Init()

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
	slog.Info("bytes read!", "n", len(bytes))

	mem := mem.NewRAM()
	mem.Copy(bytes, 0)

	cpu := cpu.NewCPU(&mem)

	for {
		cpu.FetchExecute()
		// time.Sleep(time.Second / 100)
	}
}
