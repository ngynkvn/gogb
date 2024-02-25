package cpu

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

func SplitU16(value uint16) (uint8, uint8) {
	return uint8(value >> 8), uint8(value & 0b1111_1111)
}

func (c *CPU) SetR8(reg uint8) func(uint8) {
	return func(u uint8) {
		*c.Location(reg) = u
	}
}

func (c *CPU) Set16(dst uint8, value uint16) {
	switch dst {
	case 0:
		hiDst, loDst := &c.B, &c.C
		hi, lo := SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 1:
		hiDst, loDst := &c.D, &c.E
		hi, lo := SplitU16(value)
		*hiDst = hi
		*loDst = lo
	case 2:
		hiDst, loDst := &c.H, &c.L
		hi, lo := SplitU16(value)
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

// TODO: refactor
func (c *CPU) SetHL(value uint16) {
	c.Set16(0b10, value)
}

func (c *CPU) FetchR16(reg uint8) uint16 {
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
	c.cycle += 1
	return c.ram.ReadU8(pos)
}

// Read u16 from memory, incur +2 cycle
func (c *CPU) ReadU16(pos uint16) uint16 {
	c.cycle += 2
	return c.ram.ReadU16(pos)
}

func (c *CPU) ReadU8Imm() uint8 {
	result := c.ram.ReadU8(c.PC)
	c.PC += 1
	return result
}
func (c *CPU) ReadU16Imm() uint16 {
	result := c.ram.ReadU16(c.PC)
	c.PC += 2
	return result
}

func (c *CPU) MemSet8(pos uint16) func(uint8) {
	return func(u uint8) {
		*c.ram.Ptr(pos) = u
	}
}
