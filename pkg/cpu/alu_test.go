package cpu

import (
	"gogb/pkg/mem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCp(t *testing.T) {
	mem := mem.NewRAM()
	cpu := NewCPU(mem, nil)

	cpu.A = 0x42

	cpu.InstrCp(cpu.A, cpu.A)

	assert.True(t, cpu.F_Z(), cpu.F_N())
	assert.EqualValues(t, 0x42, cpu.A)
	assert.NotEqualValues(t, 0x00, cpu.A)
}

func TestAdc(t *testing.T) {
	a := uint8(91)
	b := uint8(165)

	mem := mem.NewRAM()
	cpu := NewCPU(mem, nil)

	cpu.InstrAdd(cpu.SetA, a, b, true)

	assert.EqualValues(t, 0, cpu.A)
	assert.True(t, cpu.F_Z())
}
