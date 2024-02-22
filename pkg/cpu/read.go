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
	return uint8(value >> 8), uint8(value & 0xF)
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
	c.PC += 1
	return c.ram.ReadU8(pos)
}

// Read u16 from memory, incur +2 cycle
func (c *CPU) ReadU16(pos uint16) uint16 {
	c.cycle += 2
	c.PC += 2
	return c.ram.ReadU16(pos)
}

func (c *CPU) Ld16(opcode uint8) {
	dst := (opcode >> 3) & 0b11
	c.Set16(dst, c.ReadU16(c.PC))
}

func (c *CPU) LdMem8(opcode uint8) {
	dst := (opcode >> 4) & 0b11
	switch dst {
	case 0:
		// BC
		c.InstrLd8(c.MemSet8(c.BC()), c.A)
	case 1:
		// DE
		c.InstrLd8(c.MemSet8(c.DE()), c.A)
	case 2:
		// HL+
		c.InstrLd8(c.MemSet8(c.HL()), c.A)
	case 3:
		// HL-
		c.InstrLd8(c.MemSet8(c.HL()), c.A)
	default:
		panic("Out of range")
	}
}

func (c *CPU) MemSet8(pos uint16) func(uint8) {
	return func(u uint8) {
		*c.ram.Ptr(pos) = u
	}
}

func (c *CPU) InstrLd8(set func(uint8), value uint8) {
	set(value)
}
