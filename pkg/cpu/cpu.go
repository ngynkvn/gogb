package cpu

import (
	"fmt"
	"gogb/pkg/mem"
)

// https://gbdev.io/gb-opcodes//optables/

// https://gbdev.io/pandocs/About.html

// https://gekkio.fi/files/gb-docs/gbctr.pdf

// M-cycles
var OP_CYCLES = [256]uint8{
	////////x0,x1,x2,x3,x4,x5,x6,x7,x8,x9,xA,xB,xC,xD,xE,xF
	/* 0x */ 1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1, //
	/* 1x */ 1, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1, //
	/* 2x */ 2, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1, //
	/* 3x */ 2, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1, //
	/* 4x */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* 5x */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* 6x */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* 7x */ 2, 2, 2, 2, 2, 2, 1, 2, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* 8x */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* 9x */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* Ax */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* Bx */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, //
	/* Cx */ 2, 3, 3, 4, 3, 4, 2, 4, 2, 4, 3, 1, 3, 6, 2, 4, //
	/* Dx */ 2, 3, 3, 1, 3, 4, 2, 4, 2, 4, 3, 1, 3, 1, 2, 4, //
	/* Ex */ 3, 3, 2, 1, 1, 4, 2, 4, 4, 1, 4, 1, 1, 1, 2, 4, //
	/* Fx */ 3, 3, 2, 1, 1, 4, 2, 4, 3, 2, 4, 1, 1, 1, 2, 4, //
}

var OP_LEN = [256]uint8{
	////////x0,x1,x2,x3,x4,x5,x6,x7,x8,x9,xA,xB,xC,xD,xE,xF
	/* 0x */ 1, 3, 1, 1, 1, 1, 2, 1, 3, 1, 1, 1, 1, 1, 2, 1, //
	/* 1x */ 2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1, //
	/* 2x */ 2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1, //
	/* 3x */ 2, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 2, 1, //
	/* 4x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* 5x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* 6x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* 7x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* 8x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* 9x */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* Ax */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* Bx */ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, //
	/* Cx */ 1, 1, 3, 3, 3, 1, 2, 1, 1, 1, 3, 2, 3, 3, 2, 1, // TODO: I patched CB to be len 2 since all the prefixed ops are all len 2.
	/* Dx */ 1, 1, 3, 1, 3, 1, 2, 1, 1, 1, 3, 1, 3, 1, 2, 1, //
	/* Ex */ 2, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 1, 1, 2, 1, //
	/* Fx */ 2, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 1, 1, 2, 1, //
}

var INSTR_NAME = [256]string{
	"NOP", "LD:BC,n16", "LD:BC,A", "INC:BC", "INC:B", "DEC:B", "LD:B,n8", "RLCA", "LD:a16,SP", "ADD:HL,BC", "LD:A,BC", "DEC:BC", "INC:C", "DEC:C", "LD:C,n8", "RRCA", //
	"STOP:n8", "LD:DE,n16", "LD:DE,A", "INC:DE", "INC:D", "DEC:D", "LD:D,n8", "RLA", "JR:e8", "ADD:HL,DE", "LD:A,DE", "DEC:DE", "INC:E", "DEC:E", "LD:E,n8", "RRA", //
	"JR:NZ,e8", "LD:HL,n16", "LD:HL,A", "INC:HL", "INC:H", "DEC:H", "LD:H,n8", "DAA", "JR:Z,e8", "ADD:HL,HL", "LD:A,HL", "DEC:HL", "INC:L", "DEC:L", "LD:L,n8", "CPL", //
	"JR:NC,e8", "LD:SP,n16", "LD:HL,A", "INC:SP", "INC:HL", "DEC:HL", "LD:HL,n8", "SCF", "JR:C,e8", "ADD:HL,SP", "LD:A,HL", "DEC:SP", "INC:A", "DEC:A", "LD:A,n8", "CCF", //
	"LD:B,B", "LD:B,C", "LD:B,D", "LD:B,E", "LD:B,H", "LD:B,L", "LD:B,HL", "LD:B,A", "LD:C,B", "LD:C,C", "LD:C,D", "LD:C,E", "LD:C,H", "LD:C,L", "LD:C,HL", "LD:C,A", //
	"LD:D,B", "LD:D,C", "LD:D,D", "LD:D,E", "LD:D,H", "LD:D,L", "LD:D,HL", "LD:D,A", "LD:E,B", "LD:E,C", "LD:E,D", "LD:E,E", "LD:E,H", "LD:E,L", "LD:E,HL", "LD:E,A", //
	"LD:H,B", "LD:H,C", "LD:H,D", "LD:H,E", "LD:H,H", "LD:H,L", "LD:H,HL", "LD:H,A", "LD:L,B", "LD:L,C", "LD:L,D", "LD:L,E", "LD:L,H", "LD:L,L", "LD:L,HL", "LD:L,A", //
	"LD:HL,B", "LD:HL,C", "LD:HL,D", "LD:HL,E", "LD:HL,H", "LD:HL,L", "HALT", "LD:HL,A", "LD:A,B", "LD:A,C", "LD:A,D", "LD:A,E", "LD:A,H", "LD:A,L", "LD:A,HL", "LD:A,A", //
	"ADD:A,B", "ADD:A,C", "ADD:A,D", "ADD:A,E", "ADD:A,H", "ADD:A,L", "ADD:A,HL", "ADD:A,A", "ADC:A,B", "ADC:A,C", "ADC:A,D", "ADC:A,E", "ADC:A,H", "ADC:A,L", "ADC:A,HL", "ADC:A,A", //
	"SUB:A,B", "SUB:A,C", "SUB:A,D", "SUB:A,E", "SUB:A,H", "SUB:A,L", "SUB:A,HL", "SUB:A,A", "SBC:A,B", "SBC:A,C", "SBC:A,D", "SBC:A,E", "SBC:A,H", "SBC:A,L", "SBC:A,HL", "SBC:A,A", //
	"AND:A,B", "AND:A,C", "AND:A,D", "AND:A,E", "AND:A,H", "AND:A,L", "AND:A,HL", "AND:A,A", "XOR:A,B", "XOR:A,C", "XOR:A,D", "XOR:A,E", "XOR:A,H", "XOR:A,L", "XOR:A,HL", "XOR:A,A", //
	"OR:A,B", "OR:A,C", "OR:A,D", "OR:A,E", "OR:A,H", "OR:A,L", "OR:A,HL", "OR:A,A", "CP:A,B", "CP:A,C", "CP:A,D", "CP:A,E", "CP:A,H", "CP:A,L", "CP:A,HL", "CP:A,A", //
	"RET:NZ", "POP:BC", "JP:NZ,a16", "JP:a16", "CALL:NZ,a16", "PUSH:BC", "ADD:A,n8", "RST:$00", "RET:Z", "RET", "JP:Z,a16", "PREFIX", "CALL:Z,a16", "CALL:a16", "ADC:A,n8", "RST:$08", //
	"RET:NC", "POP:DE", "JP:NC,a16", "ILLEGAL_D3", "CALL:NC,a16", "PUSH:DE", "SUB:A,n8", "RST:$10", "RET:C", "RETI", "JP:C,a16", "ILLEGAL_DB", "CALL:C,a16", "ILLEGAL_DD", "SBC:A,n8", "RST:$18", //
	"LDH:a8,A", "POP:HL", "LD:C,A", "ILLEGAL_E3", "ILLEGAL_E4", "PUSH:HL", "AND:A,n8", "RST:$20", "ADD:SP,e8", "JP:HL", "LD:a16,A", "ILLEGAL_EB", "ILLEGAL_EC", "ILLEGAL_ED", "XOR:A,n8", "RST:$28", //
	"LDH:A,a8", "POP:AF", "LD:A,C", "DI", "ILLEGAL_F4", "PUSH:AF", "OR:A,n8", "RST:$30", "LD:HL,SP,e8", "LD:SP,HL", "LD:A,a16", "EI", "ILLEGAL_FC", "ILLEGAL_FD", "CP:A,n8", "RST:$38", //
}

type CPU struct {
	ram        *mem.RAM
	stop       bool
	halt       bool
	cycle      uint
	interrupts bool

	A, F, B, C, D, E uint8
	H, L             uint8
	SP, PC           uint16
}

func NewCPU(mem *mem.RAM) *CPU {
	cpu := CPU{
		ram: mem,
	}
	return &cpu
}

func (c *CPU) SetA(val uint8) {
	c.A = val
}

func (c *CPU) FetchExecute() {
	if c.halt {
		return
	}
	opcode := c.ReadU8Imm()
	// add todo:['0x08', '0xF2', '0xF8', '0xF9']
	switch opcode {
	case 0x00:
		// NOP
		break
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
		0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
		0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
		0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75 /*0x76*/, 0x77,
		0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F:
		c.Ld(opcode)
	case 0x01, 0x11, 0x21, 0x31:
		// LD r16, n16
		c.Ld16(opcode)
	case 0x02, 0x12, 0x22, 0x32:
		// LD [r16mem], A
		c.LdMem8(opcode)
	case 0xE2:
		// LD [C], A
		val := c.A
		pos := 0xFF00 + uint16(c.C)
		*c.ram.Ptr(pos) = val
	case 0x76:
		c.halt = true
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87:
		// ADD A, r8
		c.Add(opcode, false)
	case 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F:
		// ADC A, r8
		c.Add(opcode, true)
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
		// SUB A, r8
		c.Sub(opcode, false)
	case 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F:
		// SBC A, r8
		c.Sub(opcode, true)
	case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7:
		// AND A, r8
		c.And(opcode)
	case 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF:
		// XOR A, r8
		c.Xor(opcode)
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		// OR A, r8
		c.Or(opcode)
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		// INC A, r8
		c.Cp(opcode)
	case 0x03, 0x13, 0x23, 0x33:
		// INC r16
		c.Inc16(opcode)
	case 0x0B, 0x1B, 0x2B, 0x3B:
		// DEC r16
		c.Dec16(opcode)
	case 0x0C, 0x1C, 0x2C, 0x3C,
		0x04, 0x14, 0x24, 0x34:
		// INC r8
		c.Inc8(opcode)
	case 0x0D, 0x1D, 0x2D, 0x3D,
		0x05, 0x15, 0x25, 0x35:
		// DEC r8
		c.Dec8(opcode)
	case 0x09, 0x19, 0x29, 0x39:
		// ADD HL, r16
		c.InstrAdd16(c.SetHL, c.HL(), c.FetchR16((opcode>>4)&0b11), false)
	case 0x06, 0x16, 0x26, 0x36,
		0x0E, 0x1E, 0x2E, 0x3E:
		// LD r8, n8
		c.InstrLd8(c.SetR8((opcode>>3)&0b111), c.ReadU8Imm())
	case 0x0A, 0x1A, 0x2A, 0x3A:
		// LD A, [mem]
		c.InstrLd8(c.SetA, c.ReadU8(c.FetchR16Mem((opcode>>4)&0b11)))
	case 0x18, 0x20, 0x28, 0x30, 0x38:
		// JR
		c.Jr(opcode)
	case 0xC2, 0xC3, 0xCA, 0xD2, 0xDA, 0xE9:
		// JP
		unimplementedOp(c, opcode)
	case 0xC0, 0xC8, 0xC9, 0xD0, 0xD8:
		// RET
		unimplementedOp(c, opcode)
	case 0xC4, 0xCC, 0xCD, 0xD4, 0xDC:
		// CALL
		c.CALL(opcode)
	case 0xC1, 0xD1, 0xE1, 0xF1:
		// POP
		c.POP(opcode)
	case 0xC5, 0xD5, 0xE5, 0xF5:
		// PUSH
		c.PUSH(opcode)
	case 0xC7, 0xCF, 0xD7, 0xDF, 0xE7, 0xEF, 0xF7, 0xFF:
		// RST
		c.RST(opcode)
	case 0x07:
		// RLCA
		val := c.A
		result := (val << 1) | (val >> 7)
		c.A = result

		c.SetZ(false)
		c.SetN(false)
		c.SetH(false)
		c.SetC(val > 0x7F)
	case 0x0F:
		// RRCA
		val := c.A
		result := (val >> 1) | ((val & 1) << 7)
		c.A = result

		c.SetZ(false)
		c.SetN(false)
		c.SetH(false)
		c.SetC(result > 0x7F)
	case 0x10:
		// STOP
		c.halt = true
	case 0x17:
		// RLA
		val := c.A
		var carry uint8
		if c.F_C() {
			carry = 0b1
		}
		result := uint8(val<<1) | carry

		c.A = result
		c.SetZ(false)
		c.SetN(false)
		c.SetH(false)
		c.SetC(val > 0x7F)
	case 0x1F:
		// RRA
		val := c.A
		var carry uint8
		if c.F_C() {
			carry = 0b1000_0000
		}
		result := uint8(val>>1) | carry
		c.A = result

		c.SetZ(false)
		c.SetN(false)
		c.SetH(false)
		c.SetC((val & 1) == 1)
	case 0x27:
		// DAA
		unimplementedOp(c, opcode)
	case 0x2F:
		// CPL
		unimplementedOp(c, opcode)
	case 0x37:
		// SCF
		unimplementedOp(c, opcode)
	case 0x3F:
		// CCF
		unimplementedOp(c, opcode)
	case 0xC6:
		// ADD A, n8
		c.AddImm8(false)
	case 0xD6:
		// SUB A, n8
		c.SubImm8(false)
	case 0xE6:
		// AND A, n8
		c.AndImm8()
	case 0xF6:
		// OR A, n8
		c.OrImm8()
	case 0xCE:
		// ADC A, n8
		c.AddImm8(true)
	case 0xDE:
		// SBC A, n8
		c.SubImm8(true)
	case 0xEE:
		// XOR A, n8
		c.XorImm8()
	case 0xFE:
		// CP A, n8
		c.CpImm8()
	case 0xCB:
		// PREFIX
		c.CB(c.ReadU8Imm())
	case 0xD9:
		// RETI
		unimplementedOp(c, opcode)
	case 0xE0:
		// LDH [a8], A
		arg := c.ReadU8Imm()
		val := c.A
		pos := 0xFF00 + uint16(arg)
		*c.ram.Ptr(pos) = val
	case 0xF0:
		arg := c.A
		val := c.ReadU8Imm()
		pos := 0xFF00 + uint16(arg)
		*c.ram.Ptr(pos) = val
	case 0xF3:
		// DI
		unimplementedOp(c, opcode)
	case 0xFB:
		// EI
		unimplementedOp(c, opcode)
	case 0xD3:
		// ILLEGAL_D3
		unimplementedOp(c, opcode)
	case 0xDB:
		// ILLEGAL_DB
		unimplementedOp(c, opcode)
	case 0xDD:
		// ILLEGAL_DD
		unimplementedOp(c, opcode)
	case 0xE3:
		// ILLEGAL_E3
		unimplementedOp(c, opcode)
	case 0xE4:
		// ILLEGAL_E4
		unimplementedOp(c, opcode)
	case 0xEB:
		// ILLEGAL_EB
		unimplementedOp(c, opcode)
	case 0xEC:
		// ILLEGAL_EC
		unimplementedOp(c, opcode)
	case 0xED:
		// ILLEGAL_ED
		unimplementedOp(c, opcode)
	case 0xF4:
		// ILLEGAL_F4
		unimplementedOp(c, opcode)
	case 0xFC:
		// ILLEGAL_FC
		unimplementedOp(c, opcode)
	case 0xFD:
		// ILLEGAL_FD
		unimplementedOp(c, opcode)
	default:
		unimplementedOp(c, opcode)
	}
}

func unimplementedOp(c *CPU, opcode uint8) {
	fmt.Printf("%s\n", c.Dump())
	fmt.Printf("\n\nunimplemented:%s\t%#x\n", INSTR_NAME[opcode], opcode)
	panic("unimplemented")
}

func (c *CPU) Dump() string {
	return fmt.Sprintf(
		`AF: %04X
BC: %#04x
DE: %#04x
HL: %#04x
SP: %#04x
PC: %#04x`, c.AF(), c.BC(), c.DE(), c.HL(), c.SP, c.PC)
}
