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

func (c *CPU) Inc8(opcode uint8) {
	dst := (opcode >> 3) & 0b111
	val := c.FetchR8(dst)
	result := val + 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0xF)

	*c.Location(dst) = result
}

func (c *CPU) Dec8(opcode uint8) {
	dst := (opcode >> 3) & 0b111
	val := c.FetchR8(dst)
	result := val - 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0x0)

	*c.Location(dst) = result

}

func (c *CPU) Inc16(opcode uint8) {
	dst := (opcode >> 3) & 0b111
	val := c.FetchR16(dst)
	result := val + 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0xF)

	*c.Location(dst) = result
}

func (c *CPU) Dec16(opcode uint8) {
	dst := (opcode >> 3) & 0b111
	val := c.FetchR16(dst)
	result := val - 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0x0)

	*c.Location(dst) = result

}
