package cpu

import "gogb/pkg/mem"

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

	A, F, B, C, D, E   uint8
	HL, SP, PC         uint16
	F_Z, F_N, F_H, F_C bool
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

func (c *CPU) FetchExecute() {
	opcode := c.ram.ReadU8(c.PC)
	// LD r,r
	if opcode >= 0x40 && opcode <= 0x7F {

	}
}
