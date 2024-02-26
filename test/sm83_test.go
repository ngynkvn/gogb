package test

import (
	"encoding/json"
	c "gogb/pkg/cpu"
	"gogb/pkg/mem"
	"os"
	"path"
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

	//TODO
	IME uint8
	EI  uint8
	RAM [][]uint16 `json:"ram"`
}

type TestCase struct {
	Name    string
	Initial State
	Final   State
}

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
				t.Run(tc.Name, func(t *testing.T) {
					defer func() {
						if r := recover(); r != nil {
							t.Errorf("Failed, panicked")
						}
					}()
					// TODO: Mock ram
					ram := mem.NewRAM()
					initial := tc.Initial
					final := tc.Final
					for _, setInfo := range initial.RAM {
						pos := setInfo[0]
						value := uint8(setInfo[1])
						*ram.Ptr(pos) = value
					}
					cpu := c.NewCPU(&ram)
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

					cpu.FetchExecute()

					assert.Equal(t, final.A, cpu.A, "A")
					assert.Equal(t, final.B, cpu.B, "B")
					assert.Equal(t, final.C, cpu.C, "C")
					assert.Equal(t, final.D, cpu.D, "D")
					assert.Equal(t, final.E, cpu.E, "E")
					assert.Equal(t, final.F, cpu.F, "F")
					assert.Equal(t, final.H, cpu.H, "H")
					assert.Equal(t, final.L, cpu.L, "L")
					assert.Equal(t, final.SP, cpu.SP, "SP")
					assert.Equal(t, final.PC, cpu.PC, "PC")
				})
			}
		})
	}
}
