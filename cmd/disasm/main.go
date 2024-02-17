package main

import (
	"fmt"
	"gogb/pkg/cpu"
	"gogb/pkg/mem"
	"log/slog"
	"os"
	"strings"
)

func main() {
	slog.Info("Hello!")
	if len(os.Args) < 2 {
		slog.Info("usage: disasm [path]")
		os.Exit(1)
	}

	path := os.Args[1]
	bytes, err := os.ReadFile(path)
	if err != nil {
		slog.Error("err: %s", err)
		os.Exit(1)
	}
	slog.With("n_bytes", len(bytes)).Info("bytes read")

	ram := mem.NewRAM()

	ram.Copy(bytes, 0x0000)

	for i := 0; i < len(bytes); {
		opcode := bytes[i]
		bs := bytes[i : i+int(cpu.OP_LEN[opcode])]
		bs_fmt := []string{}
		for _, b := range bs {
			bs_fmt = append(bs_fmt, fmt.Sprintf("%02X", b))
		}
		byte_str := strings.Join(bs_fmt, " ")
		fmt.Printf("%04X: %-10s %s\n", i, byte_str, cpu.INSTR_NAME[opcode])
		i += int(cpu.OP_LEN[opcode])
	}
	slog.Info("Dump finished")
}
