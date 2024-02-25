package cpu

import (
	"fmt"
	"log/slog"
)

func (c *CPU) Jr(opcode uint8) {
	condition := (opcode >> 3) & 0b111
	next := int8(c.ReadU8(c.PC))
	var pos = int32(c.PC)
	targetAddr := int32(c.PC) + int32(next)

	// TODO: Debug Levels?
	slog.Debug(fmt.Sprintf("%#04x", pos))
	slog.Debug(fmt.Sprintf("%#02x", next))
	slog.Debug(fmt.Sprintf("%#04x", targetAddr))

	switch condition {
	case 0b011:
		// Always
		pos = targetAddr
	case 0b100:
		// NZ
		if !c.F_Z() {
			pos = targetAddr
			c.cycle++
		}
	case 0b101:
		// Z
		if c.F_Z() {
			pos = targetAddr
			c.cycle++
		}
	case 0b110:
		// NC
		if !c.F_C() {
			pos = targetAddr
			c.cycle++
		}
	case 0b111:
		// C
		if c.F_C() {
			pos = targetAddr
			c.cycle++
		}
	default:
		panic("unknown")
	}
	c.PC = uint16(pos)
}
