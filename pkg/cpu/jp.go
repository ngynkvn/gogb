package cpu

func (c *CPU) CALL(opcode uint8) {
	condition := (opcode >> 3) & 0b11
	targetAddr := c.ReadU16Imm()
	pos := c.PC
	switch {
	// Always
	case (opcode == 0b110_01_101):
		c.PushStack(pos)
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b00 && !c.F_Z():
		c.PushStack(pos)
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b01 && c.F_Z():
		c.PushStack(pos)
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b10 && !c.F_C():
		c.PushStack(pos)
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b11 && c.F_C():
		c.PushStack(pos)
		pos = targetAddr
		c.SpinCycle(1)
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
		c.SpinCycle(1)
	case condition == 0b100 && !c.F_Z():
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b101 && c.F_Z():
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b110 && !c.F_C():
		pos = targetAddr
		c.SpinCycle(1)
	case condition == 0b111 && c.F_C():
		pos = targetAddr
		c.SpinCycle(1)
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
		c.SpinCycle(1)
	case cond == 0 && !c.F_Z():
		pos = target
		c.SpinCycle(1)
	case cond == 1 && c.F_Z():
		pos = target
		c.SpinCycle(1)
	case cond == 2 && !c.F_C():
		pos = target
		c.SpinCycle(1)
	case cond == 3 && c.F_C():
		pos = target
		c.SpinCycle(1)
	}
	c.PC = uint16(pos)
}

func (c *CPU) RET(opcode uint8) {
	cond := (opcode >> 3) & 0b11
	pos := c.PC
	c.SpinCycle(1)
	switch {
	case opcode&1 == 1:
		// Always
		pos = c.PopStack()
	case cond == 0 && !c.F_Z():
		pos = c.PopStack()
		c.SpinCycle(1)
	case cond == 1 && c.F_Z():
		pos = c.PopStack()
		c.SpinCycle(1)
	case cond == 2 && !c.F_C():
		pos = c.PopStack()
		c.SpinCycle(1)
	case cond == 3 && c.F_C():
		pos = c.PopStack()
		c.SpinCycle(1)
	}
	c.PC = pos
}

func (c *CPU) RST(opcode uint8) {
	tgt := uint16((opcode>>3)&0b111) * 8
	c.PushStack(c.PC)
	c.SpinCycle(1)
	c.PC = tgt
}
