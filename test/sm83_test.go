package test

import (
	"encoding/json"
	c "gogb/pkg/cpu"
	"gogb/pkg/graphics"
	"gogb/pkg/mem"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type State struct {
	PC uint16
	SP uint16
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	F  uint8
	H  uint8
	L  uint8

	IME uint8
	EI  uint8
	RAM [][]uint16 `json:"ram"`
}

type TestCase struct {
	Name    string
	Initial State
	Final   State
	Cycles  []any
}

// Note: These tests are currently failing due to implementing the graphics unit.
// Removing the STAT updates should make them pass again
// TODO(008): Find solution to keep graphics part of testing
func TestSM83(t *testing.T) {
	filepath := "../bin/tests_sm83/v1"
	dirInfo, err := os.ReadDir(filepath)
	if err != nil {
		t.Skip("No sm83 tests found")
	}
	for _, file := range dirInfo {
		t.Run(file.Name(), func(t *testing.T) {
			var tcs []TestCase
			bytes, err := os.ReadFile(path.Join(filepath, file.Name()))
			assert.NoError(t, err)
			assert.NoError(t, json.Unmarshal(bytes, &tcs))
			for _, tc := range tcs {
				t.Run(tc.Name, func(tt *testing.T) {
					tt.Parallel()
					defer func() {
						if r := recover(); r != nil {
							t.Fatalf("Failed, panicked")
						}
					}()
					// TODO(004): Mock ram
					ram := mem.NewRAM()
					initial := tc.Initial
					final := tc.Final
					for _, setInfo := range initial.RAM {
						pos := setInfo[0]
						value := uint8(setInfo[1])
						*ram.Ptr(pos) = value
					}
					cpu := c.NewCPU(ram, graphics.NewDisplay(ram))
					cpu.A = initial.A
					cpu.B = initial.B
					cpu.C = initial.C
					cpu.D = initial.D
					cpu.E = initial.E
					cpu.F = initial.F
					cpu.H = initial.H
					cpu.L = initial.L
					cpu.PC = initial.PC
					cpu.SP = initial.SP
					cpu.IME = initial.IME != 0
					cpu.EI_QUEUED = initial.EI != 0

					before := cpu.CycleM
					cpu.Update()
					after := cpu.CycleM

					assert.Equal(tt, final.A, cpu.A, "A")
					assert.Equal(tt, final.B, cpu.B, "B")
					assert.Equal(tt, final.C, cpu.C, "C")
					assert.Equal(tt, final.D, cpu.D, "D")
					assert.Equal(tt, final.E, cpu.E, "E")
					assert.Equal(tt, final.F, cpu.F, "F")
					assert.Equal(tt, final.H, cpu.H, "H")
					assert.Equal(tt, final.L, cpu.L, "L")
					assert.Equal(tt, final.SP, cpu.SP, "SP")
					assert.Equal(tt, final.PC, cpu.PC, "PC")
					assert.Equal(tt, final.IME != 0, cpu.IME, "IME")
					opcode, err := strconv.ParseUint(tc.Name[:2], 16, 8)
					assert.NoError(tt, err)
					assert.EqualValues(tt, len(tc.Cycles), after-before, "Cycle Count: "+cpu.OpcodeName(uint8(opcode)))

					for _, setInfo := range final.RAM {
						pos := setInfo[0]
						value := uint8(setInfo[1])
						assert.EqualValues(tt, value, ram.ReadU8(pos))
						*ram.Ptr(pos) = value
					}
				})
			}
		})
	}
}
