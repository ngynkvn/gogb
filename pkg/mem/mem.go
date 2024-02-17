package mem

type RAM struct {
	memory [0x10000]byte
}

func NewRAM() RAM {
	return RAM{
		memory: [0x10000]byte{},
	}
}

func (r *RAM) Copy(bytes []byte, pos int) int {
	return copy(r.memory[pos:], bytes)
}

func (r *RAM) ReadU8(pos uint16) uint8 {
	return r.memory[pos]
}

// GB ROMs are little endian, least significant bytes come first
func (r *RAM) ReadU16(pos uint16) uint16 {
	low := r.memory[pos]
	high := r.memory[pos+1]
	return uint16(high)<<8 | uint16(low)
}
