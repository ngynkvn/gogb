package graphics

import (
	"gogb/pkg/bits"
	"gogb/pkg/mem"
)

const (
	ADDR_LCDC    = 0xFF40
	ADDR_STAT    = 0xFF41
	ADDR_LY      = 0xFF44
	ADDR_LYC     = 0xFF45
	ADDR_SCROLLY = 0xFF42
	ADDR_SCROLLX = 0xFF43
	ADDR_WINDOWY = 0xFF4A
	ADDR_WINDOWX = 0xFF4B
)

const (
	LCDMode2Bound = 456 - 80
	LCDMode3Bound = LCDMode2Bound - 172
)

type Display struct {
	ram *mem.RAM

	ScanlineCounter int
}

func NewDisplay(ram *mem.RAM) *Display {
	return &Display{
		ram: ram,
	}
}

func (d *Display) LCDC() uint8 {
	return d.ram.ReadU8(ADDR_LCDC)
}

func (d *Display) LCDEnabled() bool {
	return bits.Test(d.LCDC(), 7)
}

// Test Bit 5 and 0
func (d *Display) WindowEnabled() bool {
	return bits.Test(d.LCDC(), 5) && bits.Test(d.LCDC(), 0)
}

// Test Bit 1
func (d *Display) ObjEnabled() bool {
	return bits.Test(d.LCDC(), 1)
}
func (d *Display) WindowTileMapArea() []uint8 {
	if bits.Test(d.LCDC(), 6) {
		// 1 = $9C00-$9FFF
		return d.ram.Slice(0x9C00, 0x9FFF)
	} else {
		// 0 = $9800-$9BFF
		return d.ram.Slice(0x9800, 0x9BFF)
	}
}
func (d *Display) UnsignedAddressMode() bool {
	return bits.Test(d.LCDC(), 4)
}

func (d *Display) ObjHeight() uint8 {
	if bits.Test(d.LCDC(), 2) {
		return 2
	} else {
		return 1
	}
}
func (d *Display) BGTileMapArea() []uint8 {
	if bits.Test(d.LCDC(), 3) {
		// 1 = $9C00-$9FFF
		return d.ram.Slice(0x9C00, 0x9FFF)
	} else {
		// 0 = $9800-$9BFF
		return d.ram.Slice(0x9800, 0x9BFF)
	}
}

func (d *Display) SetLCDStatus() bool {
	status := d.ram.Ptr(ADDR_LCDC)
	ly := d.ram.Ptr(ADDR_LY)

	// TODO(006): check this routine
	if !d.LCDEnabled() {
		// d.ClearScreen()
		d.ScanlineCounter = 456
		*ly = 0
		*status &= 0b1111_1100
		return false
	}

	currentLine := *ly
	var mode uint8
	currentMode := *status & 0b11
	requestInterrupt := false
	switch {
	case currentLine >= 144:
		mode = 1
		*status = bits.Set(*status, 0)
		*status = bits.Reset(*status, 1)
		requestInterrupt = bits.Test(*status, 4)
	case d.ScanlineCounter >= LCDMode2Bound:
		mode = 2
		*status = bits.Reset(*status, 0)
		*status = bits.Set(*status, 1)
	case d.ScanlineCounter >= LCDMode3Bound:
		mode = 3
		*status = bits.Set(*status, 0)
		*status = bits.Set(*status, 1)
		if mode != currentMode {
			// Draw scanline in mode 3, this differs from actual gb
			// behavior
			d.Draw(currentLine)
		}
	default:
		mode = 0
		*status = bits.Reset(*status, 0)
		*status = bits.Reset(*status, 1)
		requestInterrupt = bits.Test(*status, 3)
		if mode != currentMode {
			// HDMA transfer
		}
	}
	if currentLine == d.ram.ReadU8(ADDR_LYC) {
		*status = bits.Set(*status, 2)
		requestInterrupt = bits.Test(*status, 6)
	} else {
		*status = bits.Reset(*status, 2)
	}
	return requestInterrupt && mode != currentMode
}

func (d *Display) Draw(scanline uint8) {
	if d.WindowEnabled() {
		d.RenderTiles(scanline)
	}
	if d.ObjEnabled() {
		d.RenderSprites(int32(scanline))
	}
}

func (d *Display) RenderTiles(scanline uint8) {

}
func (d *Display) RenderSprites(scanline int32) {

}
