package cpu

func (c *CPU) F_Z() bool {
	return (c.F & (1 << 7)) > 0
}
func (c *CPU) F_N() bool {
	return (c.F & (1 << 6)) > 0
}
func (c *CPU) F_H() bool {
	return (c.F & (1 << 5)) > 0
}
func (c *CPU) F_C() bool {
	return (c.F & (1 << 4)) > 0
}

func (c *CPU) SetZ(set bool) {
	if set {
		c.F = c.F | 1<<uint8(7)
	} else {
		c.F = c.F & ^uint8(1<<7)
	}
}
func (c *CPU) SetN(set bool) {
	if set {
		c.F = c.F | 1<<uint8(6)
	} else {
		c.F = c.F & ^uint8(1<<6)
	}
}
func (c *CPU) SetH(set bool) {
	if set {
		c.F = c.F | 1<<uint8(5)
	} else {
		c.F = c.F & ^uint8(1<<5)
	}
}
func (c *CPU) SetC(set bool) {
	if set {
		c.F = c.F | 1<<uint8(4)
	} else {
		c.F = c.F & ^uint8(1<<4)
	}
}

func (c *CPU) AF() uint16 {
	return (uint16(c.A) << 8) | (uint16(c.F))
}

func (c *CPU) BC() uint16 {
	return (uint16(c.B) << 8) | (uint16(c.C))
}

func (c *CPU) DE() uint16 {
	return (uint16(c.D) << 8) | (uint16(c.E))
}

func (c *CPU) HL() uint16 {
	return (uint16(c.H) << 8) | (uint16(c.L))
}
