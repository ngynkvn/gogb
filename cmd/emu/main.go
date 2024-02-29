package main

import (
	"flag"
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/log"
	"gogb/pkg/mem"
	"gogb/pkg/render"
	"log/slog"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	log.Init()
	flag.Parse()
	if *cpuprofile != "" {
		slog.Info("Profiling", "f", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	args := flag.Args()

	slog.Info("Hello!")
	if len(args) < 1 {
		slog.Info("usage: emu [path]")
		os.Exit(1)
	}
	mem := mem.NewRAM()
	path := args[0]
	// bootrom, err := os.ReadFile("./bootrom.bin")
	// if err != nil {
	// 	slog.Error("err: %s", err)
	// 	os.Exit(1)
	// }
	// mem.CopyBootRom(bootrom)
	// slog.Info("bytes read!", "n", len(bootrom))
	bytes, err := os.ReadFile(path)
	if err != nil {
		slog.Error("error opening rom", "err", err)
		os.Exit(1)
	}
	mem.Copy(bytes, 0)

	display := graphics.NewDisplay(mem)
	cpu := cpu.NewCPU(mem, display)
	// TODO
	cpu.SkipBootRom()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pprof.StopCPUProfile()
			slog.Info("Capture successful")
			os.Exit(0)
		}
	}()
	go func() {
		time.Sleep(2 * time.Minute)
		c <- nil
	}()

	renderer := render.NewEbiten(cpu, display, mem)
	renderer.Start()
}
