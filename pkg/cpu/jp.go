package cpu

func (c *CPU) CALL(opcode uint8) {
	condition := (opcode >> 3) & 0b11
	targetAddr := c.ReadU16Imm()
	pos := c.PC
	switch {
	// Always
	case (opcode == 0b110_01_101):
		c.SP -= 2
		c.ram.WriteU16(c.SP, pos)
		pos = targetAddr
		c.cycle++
	case condition == 0b00 && !c.F_Z():
		c.SP -= 2
		c.WriteU16(c.SP, pos)
		pos = targetAddr
		c.cycle++
	case condition == 0b01 && c.F_Z():
		c.SP -= 2
		c.WriteU16(c.SP, pos)
		pos = targetAddr
		c.cycle++
	case condition == 0b10 && !c.F_C():
		c.SP -= 2
		c.WriteU16(c.SP, pos)
		pos = targetAddr
		c.cycle++
	case condition == 0b11 && c.F_C():
		c.SP -= 2
		c.WriteU16(c.SP, pos)
		pos = targetAddr
		c.cycle++
	}
	c.PC = pos
}

func (c *CPU) Jr(opcode uint8) {
	condition := (opcode >> 3) & 0b111
	next := int8(c.ReadU8Imm())
	var pos = int32(c.PC)
	targetAddr := int32(c.PC) + int32(next)

	switch {
	case condition == 0b011:
		// Always
		pos = targetAddr
	case condition == 0b100 && !c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b101 && c.F_Z():
		pos = targetAddr
		c.cycle++
	case condition == 0b110 && !c.F_C():
		pos = targetAddr
		c.cycle++
	case condition == 0b111 && c.F_C():
		pos = targetAddr
		c.cycle++
	}
	c.PC = uint16(pos)
}

func (c *CPU) JP(opcode uint8) {
	cond := (opcode >> 3) & 0b11
	// TODO(002): refactor
	if (opcode&1) == 1 && cond == 0b01 {
		// JP HL
		c.PC = c.HL()
		return
	}

	target := c.ReadU16Imm()
	pos := c.PC
	switch {
	case (opcode & 0b111) == 0b011:
		//Always
		pos = target
		c.cycle += 1
	case cond == 0 && !c.F_Z():
		pos = target
		c.cycle += 1
	case cond == 1 && c.F_Z():
		pos = target
		c.cycle += 1
	case cond == 2 && !c.F_C():
		pos = target
		c.cycle += 1
	case cond == 3 && c.F_C():
		pos = target
		c.cycle += 1
	}
	c.PC = uint16(pos)
}

func (c *CPU) RET(opcode uint8) {
	cond := (opcode >> 3) & 0b11
	pos := c.PC
	switch {
	case opcode&1 == 1:
		// Always
		pos = c.ReadU16(c.SP)
		c.SP += 2
		c.cycle += 1
	case cond == 0 && !c.F_Z():
		pos = c.ReadU16(c.SP)
		c.SP += 2
		c.cycle += 1
	case cond == 1 && c.F_Z():
		pos = c.ReadU16(c.SP)
		c.SP += 2
		c.cycle += 1
	case cond == 2 && !c.F_C():
		pos = c.ReadU16(c.SP)
		c.SP += 2
		c.cycle += 1
	case cond == 3 && c.F_C():
		pos = c.ReadU16(c.SP)
		c.SP += 2
		c.cycle += 1
	}
	c.PC = pos
}

func (c *CPU) RETI(opcode uint8) {
}

func (c *CPU) RST(opcode uint8) {
	tgt := uint16((opcode>>3)&0b111) * 8
	c.SP -= 2
	c.ram.WriteU16(c.SP, c.PC)
	c.PC = tgt
}
