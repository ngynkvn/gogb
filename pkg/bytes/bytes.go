package bytes

func SplitU16(value uint16) (uint8, uint8) {
	return uint8(value >> 8), uint8(value & 0b1111_1111)
}
