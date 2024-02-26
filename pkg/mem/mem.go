package mem

import (
	"fmt"
	"gogb/pkg/bits"
	"log"
)

type RAM struct {
	bootrom [0x100]byte
	memory  [0x10000]byte
}

func NewRAM() RAM {
	return RAM{
		bootrom: [0x100]byte{},
		memory:  [0x10000]byte{},
	}
}

func (r *RAM) InBootRom() bool {
	return r.memory[0xFF50] == 0
}

func (r *RAM) CopyBootRom(rom []byte) {
	if len(rom) != 256 {
		log.Fatalf("Incorrect size for bootrom, expected 256 bytes but got %d\n", len(rom))
	}
	copy(r.bootrom[:], rom)
}

func (r *RAM) Copy(bytes []byte, pos int) int {
	return copy(r.memory[pos:], bytes)
}

func (r *RAM) ReadU8(pos uint16) uint8 {
	// if r.InBootRom() && pos < 0x100 {
	// 	return r.bootrom[pos]
	// }
	return r.memory[pos]
}

func (r *RAM) WriteU16(pos uint16, value uint16) {
	hi, lo := bits.SplitU16(value)
	r.memory[pos] = lo
	r.memory[pos+1] = hi
}

func (r *RAM) WriteU8(pos uint16, value uint8) {
	// TODO: Proper hook for cpu tests
	if pos == 0xFF01 {
		fmt.Printf("%c", value)
	}
	r.memory[pos] = value
}

func (r *RAM) Ptr(pos uint16) *uint8 {
	return &r.memory[pos]
}

// GB ROMs are little endian, least significant bytes come first
func (r *RAM) ReadU16(pos uint16) uint16 {
	// if r.InBootRom() && pos < 0x100 {
	// 	low := r.memory[pos]
	// 	high := r.memory[pos+1]
	// 	return uint16(high)<<8 | uint16(low)
	// }
	low := r.memory[pos]
	high := r.memory[pos+1]
	return uint16(high)<<8 | uint16(low)
}
