package cpu

func (c *CPU) Ld(opcode uint8) {
	src := opcode & 0b111
	val := c.FetchR8(src)
	dst := (opcode >> 3) & 0b111
	*c.Location(dst) = val
}
