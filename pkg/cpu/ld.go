package cpu

func (c *CPU) Ld(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	dst := (opcode >> 3) & 0b111
	*c.Location(dst) = val
}

func (c *CPU) Ld16(opcode uint8) {
	dst := (opcode >> 4) & 0b11
	c.Set16(dst, c.ReadU16Imm())
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
		c.SetHL(c.HL() + 1)
	case 3:
		// HL-
		c.InstrLd8(c.MemSet8(c.HL()), c.A)
		c.SetHL(c.HL() - 1)
	default:
		panic("Out of range")
	}
}
func (c *CPU) InstrLd8(set func(uint8), value uint8) {
	set(value)
}

func (c *CPU) PUSH(opcode uint8) {
	src := (opcode >> 4) & 0b11
	val := c.FetchR16Stk(src)
	c.SP -= 2
	c.ram.WriteU16(c.SP, val)
}
func (c *CPU) POP(opcode uint8) {
	dst := (opcode >> 4) & 0b11
	val := c.ReadU16(c.SP)
	c.Set16Stk(dst, val)
	c.SP += 2
}
