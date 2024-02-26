package cpu

import (
	"fmt"
	"log/slog"
)

// TODO: paranoid mode
func (c *CPU) CALL(opcode uint8) {
	condition := (opcode >> 3) & 0b11
	var pos = c.PC
	targetAddr := c.ReadU16Imm()
	switch {
	// Always
	case (opcode == 0b110_01_101):
		pos = targetAddr
		c.cycle++
	case condition == 0b00 && !c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b01 && c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b10 && !c.F_C():
		pos = targetAddr
		c.cycle++
	case condition == 0b11 && c.F_C():
		pos = targetAddr
		c.cycle++
	}
	c.PC = pos
}

// TODO: paranoid mode
func (c *CPU) Jr(opcode uint8) {
	condition := (opcode >> 3) & 0b111
	next := int8(c.ReadU8Imm())
	var pos = int32(c.PC)
	targetAddr := int32(c.PC) + int32(next)

	// TODO: Debug Levels?
	slog.Debug(fmt.Sprintf("%#04x", pos))
	slog.Debug(fmt.Sprintf("%#02x", next))
	slog.Debug(fmt.Sprintf("%#04x", targetAddr))

	switch {
	case condition == 0b011:
		// Always
		pos = targetAddr
	case condition == 0b100 && !c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b101 && c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b110 && !c.F_C():
		pos = targetAddr
		c.cycle++
	case condition == 0b111 && c.F_C():
		pos = targetAddr
		c.cycle++
	}
	c.PC = uint16(pos)
}

func (c *CPU) RST(opcode uint8) {
	tgt := uint16((opcode>>3)&0b11) * 8
	c.PC = tgt
}
