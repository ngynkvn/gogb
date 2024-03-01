package blargg

import (
	"fmt"
	"gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var CpuTestsExpectedCycles = map[string]uint{
	"01-special.gb":            0x2602f8,
	"03-op sp,hl.gb":           0x2629b9,
	"04-op r,imm.gb":           0x2c6429,
	"05-op rp.gb":              0x3c96d3,
	"06-ld r,r.gb":             0x9b146,
	"07-jr,jp,call,ret,rst.gb": 0xb9a3c,
	"08-misc instrs.gb":        0x88e93,
	"09-op r,r.gb":             0x92b7c8,
	"10-bit ops.gb":            0xde9f3c,
	"11-op a,(hl).gb":          0x119d172,
}

func TestCpuInstrs(t *testing.T) {
	rom_path := "cpu/"
	dirInfo, err := os.ReadDir(rom_path)

	if err != nil {
		t.Skip("Can't open directory")
	}

	for _, file := range dirInfo {
		if path.Ext(file.Name()) != ".gb" {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			bytes, err := os.ReadFile(path.Join(rom_path, file.Name()))
			assert.NoError(t, err)
			ram := mem.NewRAM()
			ram.Copy(bytes, 0)
			graphics := graphics.NewDisplay(ram)
			cpu := cpu.NewCPU(ram, graphics)
			cpu.SkipBootRom()
			for maxCycles := CpuTestsExpectedCycles[file.Name()]; maxCycles >= 0; maxCycles -= cpu.Update() {
				if strings.Contains(ram.Serial.String(), "Passed") {
					t.Log("Passed at Cycle=", fmt.Sprintf("%x", cpu.CycleM))
					return
				}
			}
			t.Fatal("Failed, cycles was", fmt.Sprintf("%x", cpu.CycleM))
		})
	}

}

func TestMemoryInstrs(t *testing.T) {
	rom_path := "memory/mem_timing.gb"
	rom, err := os.ReadFile(rom_path)
	if err != nil {
		t.Skip("Can't open directory")
	}

	ram := mem.NewRAM()
	ram.Copy(rom, 0)
	graphics := graphics.NewDisplay(ram)
	cpu := cpu.NewCPU(ram, graphics)
	cpu.SkipBootRom()
	for maxCycles := uint(0x18f6e7); maxCycles >= 0; maxCycles -= cpu.Update() {
		if strings.Contains(ram.Serial.String(), "Passed") {
			t.Log("Passed at Cycle=", fmt.Sprintf("%x", cpu.CycleM))
			return
		}
	}
	t.Fatal("Failed, cycles was", fmt.Sprintf("%x", cpu.CycleM))
}
