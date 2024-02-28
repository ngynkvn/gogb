package cpu

import (
	"gogb/pkg/bits"
	"gogb/pkg/graphics"
)

func (c *CPU) Location(reg uint8) *uint8 {
	switch reg {
	case 0:
		return &c.B
	case 1:
		return &c.C
	case 2:
		return &c.D
	case 3:
		return &c.E
	case 4:
		return &c.H
	case 5:
		return &c.L
	case 6:
		return c.ram.Ptr(c.HL())
	case 7:
		return &c.A
	default:
		panic("Out of range")
	}
}

func (c *CPU) SetR8(reg uint8) func(uint8) {
	return func(u uint8) {
		// TODO: refactor
		if reg == 6 {
			c.CycleM++
		}
		*c.Location(reg) = u
	}
}

func (c *CPU) Set16Stk(dst uint8, value uint16) {
	switch dst {
	case 0:
		hiDst, loDst := &c.B, &c.C
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 1:
		hiDst, loDst := &c.D, &c.E
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 2:
		hiDst, loDst := &c.H, &c.L
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 3:
		hi, lo := bits.SplitU16(value)
		c.A = hi
		c.F = lo & 0xF0
	default:
		panic("Out of range")
	}
}

func (c *CPU) Set16(dst uint8, value uint16) {
	switch dst {
	case 0:
		hiDst, loDst := &c.B, &c.C
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 1:
		hiDst, loDst := &c.D, &c.E
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 2:
		hiDst, loDst := &c.H, &c.L
		hi, lo := bits.SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 3:
		c.SP = value
	default:
		panic("Out of range")
	}
}

func (c *CPU) SetSP(value uint16) {
	c.Set16(0b11, value)
}

// TODO(003): refactor
func (c *CPU) SetHL(value uint16) {
	c.Set16(0b10, value)
}

func (c *CPU) FetchR16(reg uint8) uint16 {
	c.CycleM += 1
	switch reg {
	case 0:
		return c.BC()
	case 1:
		return c.DE()
	case 2:
		return c.HL()
	case 3:
		return c.SP
	default:
		panic("Out of range")
	}
}

func (c *CPU) FetchR16Stk(reg uint8) uint16 {
	switch reg {
	case 0:
		return c.BC()
	case 1:
		return c.DE()
	case 2:
		return c.HL()
	case 3:
		return c.AF()
	default:
		panic("Out of range")
	}
}

func (c *CPU) FetchR16Mem(reg uint8) uint16 {
	switch reg {
	case 0:
		return c.BC()
	case 1:
		return c.DE()
	case 2:
		// HL+
		val := c.HL()
		c.SetHL(c.HL() + 1)
		return val
	case 3:
		// HL -
		val := c.HL()
		c.SetHL(c.HL() - 1)
		return val
	default:
		panic("Out of range")
	}
}

func (c *CPU) FetchR8(reg uint8) uint8 {
	switch reg {
	case 0:
		return c.B
	case 1:
		return c.C
	case 2:
		return c.D
	case 3:
		return c.E
	case 4:
		return c.H
	case 5:
		return c.L
	case 6:
		return c.ReadU8(c.HL())
	case 7:
		return c.A
	default:
		panic("Out of range")
	}
}

// Read u8 from memory, incur +1 cycle
func (c *CPU) ReadU8(pos uint16) uint8 {
	c.CycleM += 1
	switch {
	case pos == ADDR_JOYPAD:
		return c.GetJoypadState()
	}
	return c.ram.ReadU8(pos)
}

// Read u16 from memory, incur +2 cycle
func (c *CPU) ReadU16(pos uint16) uint16 {
	c.CycleM += 2
	return c.ram.ReadU16(pos)
}

// Write u8 to memory, incur +1 cycle
func (c *CPU) WriteU8(pos uint16, value uint8) {
	switch {
	case pos == ADDR_DIV:
		c.DIV = 0x00
		c.ram.WriteU8(pos, 0x00)
	case pos == graphics.ADDR_LY:
		c.ram.WriteU8(pos, 0x00)
	case pos == graphics.ADDR_DMA:
		c.ram.DMA(value)
	case pos == 0xFF01:
		// Serial Output
		// TODO(001): Proper hook for serial output
		// fmt.Printf("%c", value)
		// r.Serial.WriteByte(value)
		c.ram.WriteU8(pos, value)
	case pos < 0x8000 || 0xFEA0 <= pos && pos < 0xFEFF:
		// Illegal rom writes
		break
	case 0xE000 <= pos && pos < 0xFE00:
		// ECHO ram
		*c.ram.Ptr(pos - 0x2000) = value
		c.ram.WriteU8(pos, value)
	default:
		c.ram.WriteU8(pos, value)
	}
	c.CycleM += 1
}

// Write u16 to memory, incur +2 cycle
func (c *CPU) WriteU16(pos uint16, value uint16) {
	c.CycleM += 2
	c.ram.WriteU16(pos, value)
}

func (c *CPU) ReadU8Imm() uint8 {
	result := c.ram.ReadU8(c.PC)
	c.CycleM += 1
	c.PC += 1
	return result
}
func (c *CPU) ReadU16Imm() uint16 {
	result := c.ram.ReadU16(c.PC)
	c.CycleM += 2
	c.PC += 2
	return result
}

func (c *CPU) MemSet8(pos uint16) func(uint8) {
	return func(u uint8) {
		c.WriteU8(pos, u)
	}
}

func (c *CPU) PushStack(address uint16) {
	c.SP -= 2
	c.WriteU16(c.SP, address)
}

func (c *CPU) PopStack() uint16 {
	result := c.ReadU16(c.SP)
	c.SP += 2
	return result
}
