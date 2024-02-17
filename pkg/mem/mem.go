package mem

type RAM struct {
	memory [0xFFFF]byte
}

func NewRAM() RAM {
	return RAM{
		memory: [0xFFFF]byte{},
	}
}

func (r *RAM) Copy(bytes []byte, pos int) int {
	return copy(r.memory[pos:], bytes)
}
