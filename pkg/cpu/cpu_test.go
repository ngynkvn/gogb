package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCpuLd(t *testing.T) {
	cpu := NewCPU(nil)
	cpu.A = 100
	// Load B, A
	cpu.Ld(0b01_000_111)
	assert.EqualValues(t, cpu.A, cpu.B)
	assert.EqualValues(t, 100, cpu.B)
	assert.EqualValues(t, 100, cpu.A)
}
