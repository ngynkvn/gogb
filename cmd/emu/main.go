package main

import (
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
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

	display := graphics.NewDisplay(mem)
	cpu := cpu.NewCPU(mem, display)
	// TODO
	cpu.SkipBootRom()

	for {
		cpu.Update()
	}
}
