package cpu

// TODO: test this
func (c *CPU) Add(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	carry := int16(0)
	if c.F_C() && addCarry {
		carry = 1
	}
	val := c.FetchR8(src)
	result := int16(c.A) + int16(val) + carry
	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((c.A&0xF)+(val&0xF)+uint8(carry) > 0xF)
	c.SetC(result > 0xFF)

	c.A = uint8(result)
}

// TODO: test this
func (c *CPU) Sub(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	carry := int16(0)
	if c.F_C() && addCarry {
		carry = 1
	}
	val := c.FetchR8(src)
	result := int16(c.A) - int16(val) - carry

	c.SetZ(result == 0)
	c.SetN(true)
	c.SetH(int16(c.A&0xF)-int16(val&0xF)-int16(carry) < 0x00)
	c.SetC(result < 0)

	c.A = uint8(result)
}

// TODO: test this
func (c *CPU) And(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	result := c.A & val

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(true)
	c.SetC(false)

	c.A = uint8(result)
}

// TODO: test this
func (c *CPU) Or(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	result := c.A | val

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(false)
	c.SetC(false)

	c.A = uint8(result)
}

// TODO: test this
func (c *CPU) Xor(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	result := c.A ^ val

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(false)
	c.SetC(false)

	c.A = uint8(result)
}

// TODO: test this
func (c *CPU) Cp(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	result := c.A - val

	c.SetZ(result == 0)
	c.SetN(true)
	c.SetH((c.A & 0xF) > (val & 0xF))
	c.SetC(c.A > val)
}
