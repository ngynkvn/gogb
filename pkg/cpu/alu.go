package cpu

// TODO: test this
func (c *CPU) Add(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	c.InstrAdd(c.SetA, c.A, c.FetchR8(src), addCarry)
}
func (c *CPU) AddImm(addCarry bool) {
	c.InstrAdd(c.SetA, c.A, c.ReadU8(c.PC), addCarry)
}

// TODO: test this
func (c *CPU) Sub(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	c.InstrSub(c.SetA, c.A, c.FetchR8(src), addCarry)
}
func (c *CPU) SubImm(addCarry bool) {
	c.InstrSub(c.SetA, c.A, c.ReadU8(c.PC), addCarry)
}

func (c *CPU) InstrAdd(set func(uint8), a uint8, b uint8, addCarry bool) {
	carry := int16(0)
	if c.F_C() && addCarry {
		carry = 1
	}
	result := int16(a) + int16(b) + carry
	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((a&0xF)+(b&0xF)+uint8(carry) > 0xF)
	c.SetC(result > 0xFF)

	set(uint8(result))
}

func (c *CPU) InstrSub(set func(uint8), a uint8, b uint8, addCarry bool) {
	carry := int16(0)
	if c.F_C() && addCarry {
		carry = 1
	}
	result := int16(a) - int16(b) - carry

	c.SetZ(result == 0)
	c.SetN(true)
	c.SetH(int16(a&0xF)-int16(b&0xF)-int16(carry) < 0x00)
	c.SetC(result < 0)

	set(uint8(result))
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

func (c *CPU) InstrAnd() {

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

	c.Set16(dst, result)
}

func (c *CPU) Dec16(opcode uint8) {
	dst := (opcode >> 3) & 0b111
	val := c.FetchR16(dst)
	result := val - 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0x0)

	c.Set16(dst, result)
}
