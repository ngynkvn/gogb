package bits

func SplitU16(value uint16) (uint8, uint8) {
	return uint8(value >> 8), uint8(value & 0b1111_1111)
}

func Test(value uint8, bit uint8) bool {
	return (value>>bit)&1 == 1
}

func Set(value uint8, bit uint8) uint8 {
	return (value) | (1 << bit)
}

func Reset(value uint8, bit uint8) uint8 {
	return (value) & ^(1 << bit)
}
