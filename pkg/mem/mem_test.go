package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemReadU16(t *testing.T) {
	// GIVEN
	mem := NewRAM()
	mem.Copy([]byte{0xFE, 0xFF}, 0)
	mem.Copy([]byte{0xFF, 0x9F}, 2)

	// WHEN
	u := mem.ReadU16(0x0000)
	u2 := mem.ReadU16(0x0002)

	// THEN
	assert.EqualValues(t, 0xFFFE, u)
	assert.EqualValues(t, 0x9FFF, u2)
}
