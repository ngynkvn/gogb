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
