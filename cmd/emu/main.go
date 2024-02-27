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
		slog.Info("usage: emu [path]")
		os.Exit(1)
	}
	mem := mem.NewRAM()
	path := os.Args[1]
	bootrom, err := os.ReadFile("./bootrom.bin")
	if err != nil {
		slog.Error("err: %s", err)
		os.Exit(1)
	}
	mem.CopyBootRom(bootrom)
	slog.Info("bytes read!", "n", len(bootrom))
	bytes, err := os.ReadFile(path)
	if err != nil {
		slog.Error("err: %s", err)
		os.Exit(1)
	}
	mem.Copy(bytes, 0)

	cpu := cpu.NewCPU(&mem)
	// TODO
	cpu.SkipBootRom()

	for {
		cpu.Update()
	}
}
