package cpu

// TODO: test this
func (c *CPU) Add(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	c.InstrAdd(c.SetA, c.A, c.FetchR8(src), addCarry)
}
func (c *CPU) AddImm8(addCarry bool) {
	c.InstrAdd(c.SetA, c.A, c.ReadU8Imm(), addCarry)
}

// TODO: test this
func (c *CPU) Sub(opcode uint8, addCarry bool) {
	src := opcode & 0b111
	c.InstrSub(c.SetA, c.A, c.FetchR8(src), addCarry)
}
func (c *CPU) SubImm8(addCarry bool) {
	c.InstrSub(c.SetA, c.A, c.ReadU8Imm(), addCarry)
}

// TODO: test this
func (c *CPU) And(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)

	c.InstrAnd(c.SetA, c.A, val)
}

func (c *CPU) AndImm8() {
	val := c.ReadU8Imm()
	c.InstrAnd(c.SetA, c.A, val)
}

func (c *CPU) Or(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)

	c.InstrOr(c.SetA, c.A, val)
}
func (c *CPU) OrImm8() {
	val := c.ReadU8Imm()
	c.InstrOr(c.SetA, c.A, val)
}

func (c *CPU) Xor(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	c.InstrXor(c.SetA, c.A, val)
}
func (c *CPU) XorImm8() {
	val := c.ReadU8Imm()
	c.InstrXor(c.SetA, c.A, val)
}

func (c *CPU) Cp(opcode uint8) {

	src := opcode & 0b111
	val := c.FetchR8(src)
	c.InstrCp(c.SetA, c.A, val)
}
func (c *CPU) CpImm8() {
	val := c.ReadU8Imm()
	c.InstrCp(c.SetA, c.A, val)
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

func (c *CPU) InstrAdd16(set func(uint16), a uint16, b uint16, addCarry bool) {
	carry := int32(B(addCarry))
	result := int32(a) + int32(b) + carry

	c.SetN(false)
	c.SetH((a&0xFF)+(b&0xFF)+uint16(carry) > 0xFF)
	c.SetC(result > 0xFFFF)

	set(uint16(result))
}

func (c *CPU) Inc16(opcode uint8) {
	dst := (opcode >> 3) & 0b11
	val := c.FetchR16(dst)
	result := val + 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0xF)

	c.Set16(dst, result)
}

func (c *CPU) Dec16(opcode uint8) {
	dst := (opcode >> 3) & 0b11
	val := c.FetchR16(dst)
	result := val - 1

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((val & 0xF) == 0x0)

	c.Set16(dst, result)
}

func (c *CPU) InstrAdd(set func(uint8), a uint8, b uint8, addCarry bool) {
	carry := int16(B(c.F_C() && addCarry))
	result := int16(a) + int16(b) + carry
	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH((a&0xF)+(b&0xF)+uint8(carry) > 0xF)
	c.SetC(result > 0xFF)

	set(uint8(result))
}

func (c *CPU) InstrSub(set func(uint8), a uint8, b uint8, addCarry bool) {
	carry := int16(B(c.F_C() && addCarry))
	result := int16(a) - int16(b) - carry

	c.SetZ(result == 0)
	c.SetN(true)
	c.SetH(int16(a&0xF)-int16(b&0xF)-int16(carry) < 0x00)
	c.SetC(result < 0)

	set(uint8(result))
}

func (c *CPU) InstrAnd(set func(uint8), a uint8, b uint8) {
	result := a & b
	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(true)
	c.SetC(false)

	set(uint8(result))
}

// TODO: test this
func (c *CPU) InstrXor(set func(uint8), a uint8, b uint8) {
	result := a ^ b

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(false)
	c.SetC(false)

	set(result)
}

// TODO: test this
func (c *CPU) InstrOr(set func(uint8), a uint8, b uint8) {
	result := a | b

	c.SetZ(result == 0)
	c.SetN(false)
	c.SetH(false)
	c.SetC(false)

	set(result)
}

// TODO: test this
func (c *CPU) InstrCp(set func(uint8), a uint8, b uint8) {
	result := a - b

	c.SetZ(result == 0)
	c.SetN(true)
	c.SetH((a & 0xF) > (b & 0xF))
	c.SetC(a > b)

	set(result)
}

func B(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
